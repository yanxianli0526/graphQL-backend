// Package orm provides `GORM` helpers for the creation, migration and access
// on the project's database
package orm

import (
	"fmt"
	config "graphql-go-template/envconfig"
	"graphql-go-template/internal/models"
	log "log"

	//Imports the database dialect of choice
	"gorm.io/driver/postgres"
	gormLogger "gorm.io/gorm/logger"

	"gorm.io/gorm"
)

// ORM struct to holds the gorm pointer to db
type GormDatabase struct {
	DB *gorm.DB
}

var Gorm = &gorm.Config{
	Logger: gormLogger.Default.LogMode(gormLogger.Error),
	// Logger: gormLogger.Default.LogMode(gormLogger.Info),
}

// Factory creates a db connection with the selected dialect and connection string
func Factory(env config.Database) (*GormDatabase, error) {
	fmt.Println("env.DBname", env.DBname)
	fmt.Println("env.DBUser", env.DBUser)
	fmt.Println("env.DBPassword", env.DBPassword)

	databaseConnect := fmt.Sprintf("sslmode=%s host=%s port=%v dbname=%s password=%s user=%s", env.DBSSLMode, env.DBHost, env.DBPort, env.DBname, env.DBPassword, env.DBUser)
	fmt.Println("databaseConnect:", databaseConnect)
	db, err := gorm.Open(postgres.Open(databaseConnect), Gorm)

	if err != nil {
		log.Panic("[ORM] err: ", err)
	}
	// Log every SQL command on dev, @prod: this should be disabled?
	// db.LogMode(true)
	// Automigrate tables
	// err = db.ServiceAutoMigration(db)

	log.Println("[ORM] Database connection initialized.")
	return &GormDatabase{DB: db}, err
}

func (d *GormDatabase) UpdateMigration() error {

	d.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	return d.DB.AutoMigrate(
		&models.Organization{},
		&models.OrganizationReceipt{},
		&models.OrganizationReceiptTemplateSetting{},
		&models.AutoTextField{},
		&models.File{},
		&models.User{},
		&models.Patient{},
		&models.BasicChargeSetting{},
		&models.NonFixedChargeRecord{},
		&models.SubsidySetting{},
		&models.OrganizationBasicChargeSetting{},
		&models.OrganizationNonFixedChargeSetting{},
		&models.TransferRefundLeave{},
		&models.PatientBill{},
		&models.PayRecordDetail{},
		&models.PayRecord{},
		// &models.FixedChargeRecord{},
		&models.DepositRecord{},
	)
}

func (d *GormDatabase) DBScript() error {
	var err error

	tx := d.DB.Session(&gorm.Session{SkipDefaultTransaction: true})
	err = tx.Transaction(func(tx *gorm.DB) error {
		// // 非固定
		// var nonFixedChargeRecords []models.NonFixedChargeRecord
		// err = d.DB.Unscoped().Find(&nonFixedChargeRecords).Error
		// if err != nil {
		// 	return err
		// }
		// var updateNonFixedChargeRecordsData []*models.NonFixedChargeRecord

		// for i := range nonFixedChargeRecords {
		// 	updateNonFixedChargeRecordData := models.NonFixedChargeRecord{
		// 		ID:      nonFixedChargeRecords[i].ID,
		// 		TaxType: nonFixedChargeRecords[i].IsTax,
		// 	}
		// 	updateNonFixedChargeRecordsData = append(updateNonFixedChargeRecordsData, &updateNonFixedChargeRecordData)
		// }

		// for i := range updateNonFixedChargeRecordsData {
		// 	err = tx.Unscoped().Model(&models.NonFixedChargeRecord{
		// 		ID: updateNonFixedChargeRecordsData[i].ID,
		// 	}).Where("tax_type is null").Update("tax_type", updateNonFixedChargeRecordsData[i].TaxType).Error
		// 	if err != nil {
		// 		return err
		// 	}
		// }

		// // 機構的非固定
		// var organizationNonFixedChargeSettings []models.OrganizationNonFixedChargeSetting
		// err = d.DB.Unscoped().Find(&organizationNonFixedChargeSettings).Error
		// if err != nil {
		// 	return err
		// }
		// var updateOrganizationNonFixedChargeSettingsData []*models.OrganizationNonFixedChargeSetting

		// for i := range organizationNonFixedChargeSettings {
		// 	updateOrganizationNonFixedChargeSettingData := models.OrganizationNonFixedChargeSetting{
		// 		ID:      organizationNonFixedChargeSettings[i].ID,
		// 		TaxType: organizationNonFixedChargeSettings[i].IsTax,
		// 	}
		// 	updateOrganizationNonFixedChargeSettingsData = append(updateOrganizationNonFixedChargeSettingsData, &updateOrganizationNonFixedChargeSettingData)
		// }

		// for i := range updateOrganizationNonFixedChargeSettingsData {
		// 	err = tx.Unscoped().Model(&models.OrganizationNonFixedChargeSetting{
		// 		ID: updateOrganizationNonFixedChargeSettingsData[i].ID,
		// 	}).Where("tax_type is null").Update("tax_type", updateOrganizationNonFixedChargeSettingsData[i].TaxType).Error
		// 	if err != nil {
		// 		return err
		// 	}
		// }

		// // 機構的固定
		// var organizationBasicChargeSettings []models.OrganizationBasicChargeSetting
		// err = tx.Unscoped().Find(&organizationBasicChargeSettings).Error
		// if err != nil {
		// 	return err
		// }
		// var updateOrganizationBasicChargeSettingsData []*models.OrganizationBasicChargeSetting

		// for i := range organizationBasicChargeSettings {
		// 	updateOrganizationBasicChargeSettingData := models.OrganizationBasicChargeSetting{
		// 		ID:      organizationBasicChargeSettings[i].ID,
		// 		TaxType: organizationBasicChargeSettings[i].IsTax,
		// 	}
		// 	updateOrganizationBasicChargeSettingsData = append(updateOrganizationBasicChargeSettingsData, &updateOrganizationBasicChargeSettingData)
		// }

		// for i := range updateOrganizationBasicChargeSettingsData {
		// 	err = d.DB.Unscoped().Model(&models.OrganizationBasicChargeSetting{
		// 		ID: updateOrganizationBasicChargeSettingsData[i].ID,
		// 	}).Where("tax_type is null").Update("tax_type", updateOrganizationBasicChargeSettingsData[i].TaxType).Error
		// 	if err != nil {
		// 		return err
		// 	}
		// }
		return nil
	})
	if err != nil {
		return err
	}
	// 把payrecord作處理

	return err
}
