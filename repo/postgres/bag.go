package repository

import (
	"context"
	"errors"
	"time"

	"gotemplate/core/domain"

	"gotemplate/logger"

	sq "github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"

	//sq "github.com/Masterminds/squirrel"

	"github.com/jackc/pgx/v5"
	//"github.com/jackc/pgx/v5/pgtype"
)

type BagRepository struct {
	Db  *DB
	log *logger.Logger
}

// NewUserRepository creates a new user repository instance
func NewBagRepository(db *DB, log *logger.Logger) *BagRepository {
	return &BagRepository{
		db,
		log,
	}
}

var BagTableColumns = struct {
	Bagid     string
	Bagname   string
	Bagweight string
}{
	Bagid:     "bag.bagid",
	Bagname:   "bag.bagname",
	Bagweight: "bag.bagweight",
}

func (br *BagRepository) GetBagByID(gctx *gin.Context, id uint64) (*domain.Bag1, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	br.log.Info("Came inside getuser by id")
	//var u1 domain.User
	//columns := generateColumnsFromStruct(u1, "select")
	query := psql.Select("*").
		From("bag").
		Where(sq.Eq{"bagid": id}).
		Limit(1)
	return SelectOne(ctx, br.Db, query, pgx.RowToAddrOfStructByName[domain.Bag1], br.log)
}

func (br *BagRepository) GetBags(gctx *gin.Context, skip, limit uint64) ([]domain.Bag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sb := psql.Select().
		Column("b.bagid AS BagID").
		Column("b.bagname AS BagName").
		Column("b.bagweight AS BagWeight").
		Column("json_agg(json_build_object('articleid', a.articleid, 'address', a.address)) AS Articles").
		Column("json_agg(json_build_object('number', p.number, 'type', p.type)) AS Phones").
		From("public.bag b").
		LeftJoin("public.articles a ON b.bagid = a.bagid").
		LeftJoin("public.user_phones p ON b.bagid = p.bagid").
		Limit(limit).
		Offset((skip - 1) * limit).
		GroupBy("b.bagid, b.bagname, b.bagweight")

	//sm := psql.Select(BagTableColumns.Bagid, BagTableColumns.Bagname, BagTableColumns.Bagweight).From("bag")
	//sb := psql.Select("*").From("bag")
	return SelectRows(ctx, br.Db, sb, pgx.RowToStructByNameLax[domain.Bag], br.log)
	// t, err := SelectRows1(ctx, br.Db, sb, pgx.RowToStructByNameLax[domain.Bag1], br.log)
	// if err != nil {
	// 	return nil, err
	// }
	// jsonData, _ := json.Marshal(t)

	// // Convert the JSON to a struct
	// var structData []domain.Bag1
	// json.Unmarshal(jsonData, &structData)

	// return structData, nil
}

func (br *BagRepository) Insertbag(gctx *gin.Context, bag domain.Bag1) (domain.Bag1, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//br.log.Debug("bagvalues:", bag.BagName, bag.BagWeight, bag.Testjson)
	query := psql.Insert("bag").Columns("bagname", "bagweight", "testjson").
		Values(bag.BagName, bag.BagWeight, bag.Testjson).Suffix("returning *")
	return InsertReturning(ctx, br.Db, query, pgx.RowToStructByName[domain.Bag1], br.log)

	//return bag,nil
}

func (br *BagRepository) InsertPiece(gctx *gin.Context, piece *domain.InternationalArticleSubpiece) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	br.log.Debug("Piece:", piece)

	query := psql.Insert("mailbooking_intl_subpiece").
		Columns("cess_amount", "cess_rate", "channeltype_cd", "compensation_cess_amount", "compensation_cess_rate", "counterno", "createdby", "createdon", "cth_cd", "ecommerce_paytranid", "ecommerce_sku", "ecommerce_url", "export_duty_amount", "export_duty_rate", "facilityid_bkg", "facilityid_upd", "hs_cd", "hs_description", "igst_amount", "igst_rate", "ipaddress_bkg", "ipaddress_upd", "mailbooking_intl_id", "mailbooking_intl_subpiece_id", "shiftno", "sp_asbl_fob_value", "sp_asbl_value_inr", "sp_comm_invoice_date", "sp_comm_invoice_no", "sp_count", "sp_inv_currency_cd", "sp_inv_currency_exchrate", "sp_invoice_lsn", "sp_invoice_value_pu", "sp_invoice_value_total", "sp_origin_currency_cd", "sp_tax_invoice_date", "sp_tax_invoice_no", "sp_unit_cd", "sp_weight_nett", "sp_weight_total", "tax_payment_channel_date", "tax_payment_channel_ref_no", "tax_payment_channel_source", "tax_payment_mode_cd", "updatedby", "updatedon", "usertype_cd").
		Values(piece.CessAmount, piece.CessRate, piece.ChannelTypeCD, piece.CompensationCessAmount, piece.CompensationCessRate, piece.CounterNo, piece.CreatedBy, piece.CreatedOn, piece.CTHCD, piece.ECommercePaytranID, piece.ECommerceSKU, piece.ECommerceURL, piece.ExportDutyAmount, piece.ExportDutyRate, piece.FacilityIDBKG, piece.FacilityIDUPD, piece.HSCD, piece.HSDescription, piece.IGSTAmount, piece.IGSTRate, piece.IPAddressBKG, piece.IPAddressUPD, piece.MailBookingIntlID, piece.ID, piece.ShiftNo, piece.SPAsblFOBValue, piece.SPAsblValueINR, piece.SPCommInvoiceDate, piece.SPCommInvoiceNo, piece.SPCount, piece.SPInvCurrencyCD, piece.SPInvCurrencyExchrate, piece.SPInvoiceLSN, piece.SPInvoiceValuePU, piece.SPInvoiceValueTotal, piece.SPOriginCurrencyCD, piece.SPTaxInvoiceDate, piece.SPTaxInvoiceNo, piece.SPUnitCD, piece.SPWeightNett, piece.SPWeightTotal, piece.TaxPaymentChannelDate, piece.TaxPaymentChannelRefNo, piece.TaxPaymentChannelSource, piece.TaxPaymentModeCD, piece.UpdatedBy, piece.UpdatedOn, piece.UserTypeCD)
	//query := psql.Insert("mailbooking_intl_subpiece").SetMap(generateMapFromStruct(piece, "json"))
	p, err := Insert(ctx, br.Db, query, br.log)
	br.log.Debug(p)
	return err
	//return Insert(ctx, ur.Db, query, pgx.RowToAddrOfStructByPos[domain.UserDB], ur.log)
}

type subpiecearray struct {
	p []domain.InternationalArticleSubpiece
}

// type ResultCollector[T any] struct {
// 	results []T
// 	err     error
// }

// func NewResultCollector[T any]() *ResultCollector[T] {
// 	return &ResultCollector[T]{results: make([]T, 0)}
// }

type Combinedstruct struct {
	subp  domain.InternationalArticleSubpiece
	asubp []domain.InternationalArticleSubpiece
}

type Combinedstruct1 struct {
	Subp  domain.InternationalArticleSubpiecedb
	Asubp []domain.InternationalArticleSubpiecedb
}

// func (br *BagRepository) Updatepieceswithreturn(gctx *gin.Context, intlSubpieces domain.ISubpieces) (Combinedstruct, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
// 	defer cancel()
// 	var asubp []domain.InternationalArticleSubpiece
// 	var subp domain.InternationalArticleSubpiece
// 	//collector := NewResultCollector[domain.InternationalArticleSubpiece]()
// 	errTx := br.Db.WithTx(ctx, func(tx pgx.Tx) error {

// 		batch := &pgx.Batch{}

// 		queryupdatewoSubpiece := psql.Update("mailbooking_intl_subpiece").
// 			Set("sp_unit_cd", "test23").
// 			Where(sq.Eq{"mailbooking_intl_id": 0}).
// 			Suffix("RETURNING *")

// 		// sql, args, err := queryupdatewoSubpiece.ToSql()
// 		// if err != nil {
// 		// 	return err
// 		// }

// 		// batch.Queue(sql, args...)
// 		QueueReturnRow(batch, queryupdatewoSubpiece, pgx.RowToStructByName[domain.InternationalArticleSubpiece], br.log, &subp)

// 		for _, subpiece := range intlSubpieces.IntlSubpieces {
// 			intlSubPieceSetToMap := StructToSetMap(&subpiece)

// 			updateBuilder := psql.Update("mailbooking_intl_subpiece").
// 				SetMap(intlSubPieceSetToMap).
// 				Where(sq.And{
// 					sq.Eq{"mailbooking_intl_subpiece_id": subpiece.ID},
// 					sq.Eq{"mailbooking_intl_id": subpiece.MailBookingIntlID},
// 				})

// 			sql, args, err := updateBuilder.ToSql()
// 			if err != nil {
// 				return err
// 			}

// 			batch.Queue(sql, args...)
// 			//QueueExecRow(batch, sql, args)
// 		}

// 		// queryselectSubpiece := psql.Select("*").
// 		// 	From("mailbooking_intl_subpiece").
// 		// 	Where(sq.Eq{"mailbooking_intl_id": int64(1)})

// 		// sql1, args1, err := queryselectSubpiece.ToSql()
// 		// if err != nil {
// 		// 	return err
// 		// }
// 		// batch.Queue(sql1, args1)

// 		// QueueReturn(batch, sql1, args1, pgx.RowToStructByName[domain.InternationalArticleSubpiece], &asubp)
// 		err1 := tx.SendBatch(ctx, batch).Close()
// 		if err1 != nil {
// 			return err1
// 		}

// 		return nil
// 	})

// 	if errTx != nil {
// 		return Combinedstruct{}, errTx
// 	}

// 	return Combinedstruct{subp, asubp}, nil

// }

func (br *BagRepository) UpdatepieceswithTransaction(gctx *gin.Context, intlSubpieces domain.ISubpieces) (Combinedstruct1, error) {
	//var id1 int
	var id2 int64

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var arp []domain.InternationalArticleSubpiecedb
	var a domain.InternationalArticleSubpiecedb

	TxDB := br.Db.WithTx(ctx, func(tx pgx.Tx) error {

		for _, subpiece := range intlSubpieces.IntlSubpieces {

			intlSubPieceSetToMap := StructToSetMap(&subpiece)

			updateBuilder := psql.Update("mailbooking_intl_subpiece").
				SetMap(intlSubPieceSetToMap).
				Where(sq.And{
					sq.Eq{"mailbooking_intl_subpiece_id": subpiece.ID},
					sq.Eq{"mailbooking_intl_id": subpiece.MailBookingIntlID},
				})

			id2 = subpiece.MailBookingIntlID

			err := TxExec(ctx, tx, updateBuilder, br.log)
			if err != nil {
				return err
			}

		}

		queryupdatewoSubpiece := psql.Update("mailbooking_intl_subpiece").
			Set("sp_unit_cd", "test23").
			Where(sq.Eq{"mailbooking_intl_id": id2}).
			Suffix("RETURNING *")
		err := TxReturnRow(ctx, tx, queryupdatewoSubpiece, pgx.RowToStructByName[domain.InternationalArticleSubpiecedb], br.log, &a)
		if err != nil {
			return err
		}

		queryselectSubpiece := psql.Select(" * ").
			From("mailbooking_intl_subpiece").
			Where(sq.Eq{"mailbooking_intl_id": id2})

		err = TxRows(ctx, tx, queryselectSubpiece, pgx.RowToStructByName[domain.InternationalArticleSubpiecedb], br.log, &arp)

		if err != nil {
			return err
		}

		return nil

	})
	//},pgx.Serializable)

	if TxDB != nil {
		return Combinedstruct1{}, TxDB
	}
	c := Combinedstruct1{a, arp}
	return c, nil

}

func (br *BagRepository) Updatepieceswithbatch(gctx *gin.Context, intlSubpieces domain.ISubpieces) (Combinedstruct1, error) {
	var id2 int64
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	batch := &pgx.Batch{}
	var arp []domain.InternationalArticleSubpiecedb
	var a domain.InternationalArticleSubpiecedb

	for _, subpiece := range intlSubpieces.IntlSubpieces {
		intlSubPieceSetToMap := StructToSetMap(&subpiece)

		updateBuilder := psql.Update("mailbooking_intl_subpiece").
			SetMap(intlSubPieceSetToMap).
			Where(sq.And{
				sq.Eq{"mailbooking_intl_subpiece_id": subpiece.ID},
				sq.Eq{"mailbooking_intl_id": subpiece.MailBookingIntlID},
			})
			//just for the sake of example took ID from here and initialised to id2
		id2 = subpiece.MailBookingIntlID
		QueueExecRow(batch, updateBuilder, br.log)

	}
	id1 := id2 - 1
	//id1 is one number less than actual so that we update different row than the actual passed one.
	//This is just to showcase returnrow in template instead of return as return is implemented in select.
	queryupdatewoSubpiece := psql.Update("mailbooking_intl_subpiece").
		Set("sp_unit_cd", "test23").
		Where(sq.Eq{"mailbooking_intl_id": id1}).
		Suffix("RETURNING *")

	QueueReturnRow(batch, queryupdatewoSubpiece, pgx.RowToStructByName[domain.InternationalArticleSubpiecedb], br.log, &a)

	queryselectSubpiece := psql.Select(" * ").
		From("mailbooking_intl_subpiece").
		Where(sq.Eq{"mailbooking_intl_id": id2})
	QueueReturn(batch, queryselectSubpiece, pgx.RowToStructByName[domain.InternationalArticleSubpiecedb], br.log, &arp)

	results := br.Db.SendBatch(ctx, batch).Close()

	if results != nil {
		br.log.Debug("Error results:", results)
		return Combinedstruct1{}, results
	}
	c := Combinedstruct1{a, arp}
	return c, nil
}

func (br *BagRepository) Insertbagswithsquirrel(gctx *gin.Context, bags domain.Bags) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	insertBuilder := psql.Insert("bag").
		Columns("bagname", "bagweight")
	for _, bag := range bags.Bags {
		insertBuilder = insertBuilder.Values(bag.BagName, bag.BagWeight)
	}
	confirm, err := Insert(ctx, br.Db, insertBuilder, br.log)
	br.log.Debug("Inserted bags with squirrel", confirm)

	if err != nil {
		return err
	}
	return nil

}

func (br *BagRepository) Insertbagswithpgx(gctx *gin.Context, bags domain.Bags) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	copycount, err := br.Db.CopyFrom(
		ctx,
		pgx.Identifier{"bag"},
		[]string{"bagname", "bagweight"},
		pgx.CopyFromSlice(len(bags.Bags), func(i int) ([]interface{}, error) {
			return []interface{}{bags.Bags[i].BagName, bags.Bags[i].BagWeight}, nil
		}))

	br.log.Debug("Copy Count", copycount)

	if err != nil {
		br.log.Debug("Error inserting bags:", err)
		return err
	}

	return nil

}

// func (br *BagRepository) Inserttx(gctx *gin.Context, bag domain.Bag, articles domain.Article) error {
// 	err := Tx(ctx, br.Db, br.InsertBagArticle, bag, articles)
// 	if err != nil {
// 		br.log.Debug("Error executing", err)
// 		return err
// 	}

// 	return nil

// }

func (br *BagRepository) InsertBagArticle(ctx context.Context, gctx *gin.Context, tx pgx.Tx, params ...interface{}) error {
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	var articles domain.Article
	var bag domain.Bag

	for _, param := range params {
		switch p := param.(type) {
		case domain.Article:
			articles = p
		case domain.Bag:
			bag = p
		default:
			return errors.New("unsupported parameter type")

		}
	}

	insertBagBuilder := psql.Insert("bag").Columns("bagname", "bagweight").Values(bag.BagName, bag.BagWeight)
	insertBagQuery, insertBagArgs, err := insertBagBuilder.ToSql()
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, insertBagQuery, insertBagArgs...)
	if err != nil {
		return err
	}

	insertArticleBuilder := psql.Insert("articles").Columns("address").Values(articles.Address)
	insertArticleQuery, insertArticleArgs, err := insertArticleBuilder.ToSql()
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, insertArticleQuery, insertArticleArgs...)
	if err != nil {
		return err
	}
	return nil
}

func (br *BagRepository) InsertDataBulk(ctx context.Context, gctx *gin.Context, tx pgx.Tx, params ...interface{}) error {
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	var bags []domain.Bag1
	var articles []domain.Article
	var phones []domain.Phone
	var err error

	for _, param := range params {
		switch p := param.(type) {
		case []domain.Bag1:
			bags = p
		case []domain.Article:
			articles = p

		case []domain.Phone:
			phones = p
		default:
			return errors.New("unsupported parameter type")

		}
	}

	// Insert bags
	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"bag"},
		[]string{"bagname", "bagweight"},
		pgx.CopyFromSlice(len(bags), func(i int) ([]interface{}, error) {
			return []interface{}{bags[i].BagName, bags[i].BagWeight}, nil
		}),
	)
	if err != nil {
		br.log.Debug("Error inserting bags in bulk:", err)
		return err
	}

	// Insert articles
	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"articles"},
		[]string{"address"},
		pgx.CopyFromSlice(len(articles), func(i int) ([]interface{}, error) {
			return []interface{}{articles[i].Address}, nil
		}),
	)
	if err != nil {
		br.log.Debug("Error inserting articles in bulk:", err)
		return err
	}

	// Insert phones
	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"user_phones"},
		[]string{"number", "type"},
		pgx.CopyFromSlice(len(phones), func(i int) ([]interface{}, error) {
			return []interface{}{phones[i].Number, phones[i].Type}, nil
		}),
	)
	if err != nil {
		br.log.Debug("Error inserting phones in bulk:", err)
		return err
	}

	return nil

}
