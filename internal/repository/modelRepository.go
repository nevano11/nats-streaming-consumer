package repository

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"nats-streaming-consumer/internal/entity"
	"strconv"
	"strings"
	"time"
)

type PostgresModelRepository struct {
	db *sqlx.DB
}

func NewPostgresModelRepository(db *sqlx.DB) *PostgresModelRepository {
	return &PostgresModelRepository{
		db: db,
	}
}

// ADD NEW MODEL
func (s *PostgresModelRepository) AddNewModel(model entity.Model) (int, error) {
	logrus.Debugf("Save model %s on db", model.String())
	// add delivery
	deliveryId, err := s.saveDelivery(model.Delivery)
	if err != nil {
		logrus.Warningf("Failed to add delivery: %s", err.Error())
		return 0, err
	}

	// add payment
	paymentId, err := s.savePayment(model.Payment)
	if err != nil {
		logrus.Warningf("Failed to add payment: %s", err.Error())
		return 0, err
	}

	// add model
	modelId, err := s.saveModel(model, deliveryId, paymentId)
	if err != nil {
		logrus.Warningf("Failed to add model: %s", err.Error())
		return 0, err
	}

	// add items
	s.saveItems(model.Items, modelId)

	logrus.Debugf("Model created with id = %d", modelId)
	return modelId, nil
}

func (s *PostgresModelRepository) saveDelivery(delivery entity.Delivery) (int, error) {
	deliveryId := 0
	row := s.db.QueryRow(
		"insert into delivery (name, phone, zip, city, address, region, email) values ($1, $2, $3, $4, $5, $6, $7) returning id",
		delivery.Name,
		delivery.Phone,
		delivery.Zip,
		delivery.City,
		delivery.Address,
		delivery.Region,
		delivery.Email)

	err := row.Scan(&deliveryId)
	if err != nil {
		return deliveryId, err
	}
	return deliveryId, nil
}

func (s *PostgresModelRepository) savePayment(payment entity.Payment) (int, error) {
	paymentId := 0
	row := s.db.QueryRow(
		"insert into payment (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id",
		payment.Transaction,
		payment.RequestId,
		payment.Currency,
		payment.Provider,
		payment.Amount,
		payment.PaymentDt,
		payment.Bank,
		payment.DeliveryCost,
		payment.GoodsTotal,
		payment.CustomFee)

	err := row.Scan(&paymentId)
	if err != nil {
		return paymentId, err
	}
	return paymentId, nil
}

func (s *PostgresModelRepository) saveModel(model entity.Model, deliveryId, paymentId int) (int, error) {
	modelId := 0
	row := s.db.QueryRow(
		"insert into model (order_uid, track_number, entry, delivery_id, payment_id, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) returning id",
		model.OrderUid,
		model.TrackNumber,
		model.Entry,
		deliveryId,
		paymentId,
		model.Locale,
		model.InternalSignature,
		model.CustomerId,
		model.DeliveryService,
		model.Shardkey,
		model.SmId,
		model.DateCreated,
		model.OofShard)

	err := row.Scan(&modelId)
	if err != nil {
		return modelId, err
	}
	return modelId, nil
}

func (s *PostgresModelRepository) saveItems(items []entity.Item, modelId int) {
	query := "insert into item (model_id, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) values "
	queryValuesPart := make([]string, len(items), len(items))
	parametersOnItemCount := 12
	queryParameters := make([]interface{}, len(items)*parametersOnItemCount, len(items)*parametersOnItemCount) // 12 parameters

	j := 0
	queryParamMini := make([]string, parametersOnItemCount, parametersOnItemCount)
	for i, v := range items {
		for k := 0; k < parametersOnItemCount; k++ {
			queryParamMini[k] = "$" + strconv.Itoa(i*parametersOnItemCount+k+1)
		}
		queryValuesPart[i] = "(" + strings.Join(queryParamMini, ",") + ")"
		queryParameters[j*parametersOnItemCount] = modelId
		queryParameters[j*parametersOnItemCount+1] = v.ChrtId
		queryParameters[j*parametersOnItemCount+2] = v.TrackNumber
		queryParameters[j*parametersOnItemCount+3] = v.Price
		queryParameters[j*parametersOnItemCount+4] = v.Rid
		queryParameters[j*parametersOnItemCount+5] = v.Name
		queryParameters[j*parametersOnItemCount+6] = v.Sale
		queryParameters[j*parametersOnItemCount+7] = v.Size
		queryParameters[j*parametersOnItemCount+8] = v.TotalPrice
		queryParameters[j*parametersOnItemCount+9] = v.NmId
		queryParameters[j*parametersOnItemCount+10] = v.Brand
		queryParameters[j*parametersOnItemCount+11] = v.Status
		j++
	}
	query = query + strings.Join(queryValuesPart, ",")
	s.db.MustExec(query, queryParameters...)
}

// SELECT MODEL BY uid
func (s *PostgresModelRepository) SelectModelByUid(uid string) (entity.Model, error) {
	query := `SELECT m.id as id,order_uid,m.track_number as track_number,entry,locale,internal_signature,customer_id,delivery_service,shardkey,sm_id,date_created,oof_shard,d.name delivery_name,phone,zip,city,address,region,email,transaction,request_id,currency,provider,amount,payment_dt,bank,delivery_cost,goods_total,custom_fee,chrt_id,i.track_number as item_track_number,price,rid,i.name as item_name,sale,size,total_price,nm_id,brand,status 
	FROM model m 
	JOIN delivery d on m.delivery_id = d.id 
    JOIN payment p on m.payment_id = p.id 
    LEFT JOIN item i on m.id = i.model_id 
    WHERE m.order_uid = $1`

	rows, err := s.db.Query(query, uid)
	if err != nil {
		return entity.Model{}, err
	}

	model := entity.Model{}
	isEmpty := true
	for rows.Next() {
		rowModel, _, err := s.readModel(rows)
		if err != nil {
			return entity.Model{}, err
		}
		if isEmpty {
			model = rowModel
			isEmpty = false
		} else {
			model.Items = append(model.Items, rowModel.Items...)
		}
	}
	err = rows.Err()
	if err != nil {
		return entity.Model{}, err
	}

	return model, err
}

// SELECT ALL MODELS
func (s *PostgresModelRepository) SelectAllModels() ([]entity.Model, error) {
	query := `SELECT m.id as id,order_uid,m.track_number as track_number,entry,locale,internal_signature,customer_id,delivery_service,shardkey,sm_id,date_created,oof_shard,d.name delivery_name,phone,zip,city,address,region,email,transaction,request_id,currency,provider,amount,payment_dt,bank,delivery_cost,goods_total,custom_fee,chrt_id,i.track_number as item_track_number,price,rid,i.name as item_name,sale,size,total_price,nm_id,brand,status 
	FROM model m 
	JOIN delivery d on m.delivery_id = d.id 
    JOIN payment p on m.payment_id = p.id 
    LEFT JOIN item i on m.id = i.model_id 
    ORDER BY m.id`

	rows, err := s.db.Query(query)
	if err != nil {
		return []entity.Model{}, err
	}

	resultMap := make(map[int]entity.Model)
	for rows.Next() {
		rowModel, modelId, err := s.readModel(rows)
		if err != nil {
			return []entity.Model{}, err
		}
		if val, isExists := resultMap[modelId]; isExists {
			// Add new item on model
			val.Items = append(val.Items, rowModel.Items...)
			resultMap[modelId] = val
		} else {
			// New model
			resultMap[modelId] = rowModel
		}
	}
	err = rows.Err()
	if err != nil {
		return []entity.Model{}, err
	}

	result := make([]entity.Model, 0)
	for _, v := range resultMap {
		result = append(result, v)
	}

	return result, err
}

func (s *PostgresModelRepository) readModel(row *sql.Rows) (entity.Model, int, error) {
	var id int
	var order_uid string
	var track_number string
	var entry string
	var locale string
	var internal_signature string
	var customer_id string
	var delivery_service string
	var shardkey string
	var sm_id int
	var date_created time.Time
	var oof_shard string
	var delivery_name string
	var phone string
	var zip string
	var city string
	var address string
	var region string
	var email string
	var transaction string
	var request_id string
	var currency string
	var provider string
	var amount int
	var payment_dt int
	var bank string
	var delivery_cost int
	var goods_total int
	var custom_fee int
	var chrt_id sql.NullInt32
	var item_track_number sql.NullString
	var price sql.NullInt32
	var rid sql.NullString
	var item_name sql.NullString
	var sale sql.NullInt32
	var size sql.NullString
	var total_price sql.NullInt32
	var nm_id sql.NullInt32
	var brand sql.NullString
	var status sql.NullInt32

	err := row.Scan(&id, &order_uid, &track_number, &entry, &locale, &internal_signature, &customer_id, &delivery_service, &shardkey, &sm_id, &date_created, &oof_shard, &delivery_name, &phone, &zip, &city, &address, &region, &email, &transaction, &request_id, &currency, &provider, &amount, &payment_dt, &bank, &delivery_cost, &goods_total, &custom_fee, &chrt_id, &item_track_number, &price, &rid, &item_name, &sale, &size, &total_price, &nm_id, &brand, &status)
	if err != nil {
		return entity.Model{}, 0, err
	}

	items := []entity.Item{}
	if chrt_id.Valid {
		items = append(items, entity.Item{
			ChrtId:      int(chrt_id.Int32),
			TrackNumber: item_track_number.String,
			Price:       int(price.Int32),
			Rid:         rid.String,
			Name:        item_name.String,
			Sale:        int(sale.Int32),
			Size:        size.String,
			TotalPrice:  int(total_price.Int32),
			NmId:        int(nm_id.Int32),
			Brand:       brand.String,
			Status:      int(status.Int32),
		})
	}

	return entity.Model{
		OrderUid:    order_uid,
		TrackNumber: track_number,
		Entry:       entry,
		Delivery: entity.Delivery{
			Name:    delivery_name,
			Phone:   phone,
			Zip:     zip,
			City:    city,
			Address: address,
			Region:  region,
			Email:   email,
		},
		Payment: entity.Payment{
			Transaction:  transaction,
			RequestId:    request_id,
			Currency:     currency,
			Provider:     provider,
			Amount:       amount,
			PaymentDt:    payment_dt,
			Bank:         bank,
			DeliveryCost: delivery_cost,
			GoodsTotal:   goods_total,
			CustomFee:    custom_fee,
		},
		Items:             items,
		Locale:            locale,
		InternalSignature: internal_signature,
		CustomerId:        customer_id,
		DeliveryService:   delivery_service,
		Shardkey:          shardkey,
		SmId:              sm_id,
		DateCreated:       date_created,
		OofShard:          oof_shard,
	}, id, nil
}
