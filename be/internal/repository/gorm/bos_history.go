package gorm

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"text/template"
	"time"

	"github.com/ramadhia/bosnet/be/internal/model"
	"github.com/ramadhia/bosnet/be/internal/repository"

	"gorm.io/gorm"
)

type HistoryGorm struct {
	db *gorm.DB
}

func NewHistoryRepo(db *gorm.DB) *HistoryGorm {
	return &HistoryGorm{
		db: db,
	}
}

func (p *HistoryGorm) FetchHistory(ctx context.Context, filter repository.FetchHistoryFilter) (ret []*model.History, err error) {
	q := p.db.WithContext(ctx)

	if filter.Limit != nil && filter.Offset != nil {
		limit := 5
		if filter.Limit != nil {
			limit = *filter.Limit
		}

		offset := 0
		if filter.Offset != nil {
			offset = *filter.Offset
		}

		if filter.SzAccountId != nil {
			q = q.Where("szAccountId = ?", filter.SzAccountId)
		}

		// filter by range date
		if v := filter.StartDtmTransaction; v != nil {
			q = q.Where("dtmTransaction >= ?", *v)
		}

		if v := filter.EndDtmTransaction; v != nil {
			q = q.Where("dtmTransaction <= ?", *v)
		}

		q = q.Limit(limit).Offset(offset)
	}

	var items []BosHistory
	if err = q.Find(&items).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	fmt.Println(items)

	return BosHistory{}.ToModels(items), nil
}

func (p *HistoryGorm) AddHistory(ctx context.Context, args repository.AddHistoryFilter) (bool, error) {
	q := p.db.WithContext(ctx)

	// Eksekusi transaksi dengan Snapshot Isolation
	tmpl, err := template.New("query").Parse(`
		SET TRANSACTION ISOLATION LEVEL SNAPSHOT;
		BEGIN TRANSACTION;

		DECLARE @lastNumber INT;
		DECLARE @newLastNumber INT;
		DECLARE @szTransactionId NVARCHAR(50);

		-- Ambil nilai lastnumber dari tabel BOS_Counter
		SELECT @lastNumber = iLastNumber FROM BOS_Counter WHERE szCounterId = '001-COU';

		-- Buat history ID dan hitung nilai last number yang baru
		SET @newLastNumber = @lastNumber + 1;
		-- SET @szTransactionId = FORMAT(GETDATE(), 'yyyyMMdd') + '-' + RIGHT('00000' + CAST(@newLastNumber AS VARCHAR(5)), 5) + '.0000' + CAST(@newLastNumber AS VARCHAR(5));
		SET @szTransactionId = FORMAT(GETDATE(), 'yyyyMMdd') + '-' + RIGHT('00000' + CAST(0 AS VARCHAR(5)), 5) + '.' + RIGHT('00000' + CAST(@newLastNumber AS VARCHAR(5)), 5);


		-- Insert data ke tabel BOS_History
		IF '{{.SzNote}}' = 'TRANSFER'
		BEGIN
			-- INSERT WHEN TRANSFER
			INSERT INTO BOS_History VALUES 
				(@szTransactionId, '{{.FromSzAccountId}}', '{{.SzCurrencyId}}', GETDATE(), '-{{.DecAmount}}', '{{.SzNote}}'),
				(@szTransactionId, '{{.ToSzAccountId}}', '{{.SzCurrencyId}}', GETDATE(), '{{.DecAmount}}', '{{.SzNote}}');

			-- UPDATE Table Balance
			UPDATE BOS_Balance SET decAmount = decAmount - {{.DecAmount}} WHERE szAccountId = '{{.FromSzAccountId}}' AND szCurrencyId = '{{.SzCurrencyId}}';
			UPDATE BOS_Balance SET decAmount = decAmount + {{.DecAmount}}  WHERE szAccountId = '{{.ToSzAccountId}}' AND szCurrencyId = '{{.SzCurrencyId}}';
		END
		ELSE
		BEGIN
			INSERT INTO BOS_History VALUES (@szTransactionId, '{{.ToSzAccountId}}', '{{.SzCurrencyId}}', GETDATE(), '{{.DecAmount}}', '{{.SzNote}}');
			
			-- UPDATE Table Balance
			IF '{{.SzNote}}' = 'TARIK'
			BEGIN
				UPDATE BOS_Balance SET decAmount = decAmount + {{.DecAmount}} WHERE szAccountId = '{{.ToSzAccountId}}' AND szCurrencyId = '{{.SzCurrencyId}}';
			END
			ELSE
			BEGIN
				UPDATE BOS_Balance SET decAmount = decAmount + {{.DecAmount}}  WHERE szAccountId = '{{.ToSzAccountId}}' AND szCurrencyId = '{{.SzCurrencyId}}';
			END
		END

		-- Update lastnumber di tabel BOS_Counter
		UPDATE BOS_Counter SET iLastNumber = @newLastNumber WHERE szCounterId = '001-COU';

		COMMIT TRANSACTION;`)
	if err != nil {
		panic(err)
	}

	var rawQuery bytes.Buffer
	if err := tmpl.Execute(&rawQuery, args); err != nil {
		panic(err)
	}

	if err := q.Exec(rawQuery.String()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

type BosHistory struct {
	SzTransactionId *string    `json:"szTransactionId" gorm:"column:szTransactionId"`
	SzAccountId     *string    `gorm:"column:szAccountId"`
	SzCurrencyId    *string    `gorm:"column:szCurrencyId"`
	DecAmount       *float32   `gorm:"column:decAmount"`
	SzNote          *string    `gorm:"column:szNote"`
	DtmTransaction  *time.Time `gorm:"column:dtmTransaction"`
}

func (BosHistory) FromModel(model model.History) *BosHistory {
	return &BosHistory{
		SzTransactionId: model.SzTransactionId,
		SzAccountId:     model.SzAccountId,
		SzCurrencyId:    model.SzCurrencyId,
		DecAmount:       model.DecAmount,
		SzNote:          model.SzNote,
		DtmTransaction:  model.DtmTransaction,
	}
}

func (BosHistory) TableName() string {
	return "BOS_History"
}

func (b BosHistory) ToModel() *model.History {

	//var tblCctv *model.TblCctv
	//if t.TblCctv != nil {
	//	tblCctv = t.TblCctv.ToModel()
	//
	//}

	m := &model.History{
		SzTransactionId: b.SzTransactionId,
		SzAccountId:     b.SzAccountId,
		SzCurrencyId:    b.SzCurrencyId,
		DecAmount:       b.DecAmount,
		SzNote:          b.SzNote,
		DtmTransaction:  b.DtmTransaction,
	}

	return m
}

func (b BosHistory) ToModels(d []BosHistory) []*model.History {
	fmt.Println(d)
	var models []*model.History
	for _, v := range d {
		models = append(models, v.ToModel())
	}
	return models
}
