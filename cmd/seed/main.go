package main

import (
    "log/slog"
    "math"
    "time"

    "github.com/IslamCHup/coworking-manager-project/internal/config"
    "github.com/IslamCHup/coworking-manager-project/internal/models"
    "gorm.io/gorm"
)

func main() {
    logger := config.InitLogger()
    db := config.SetupDataBase(logger)

    autoMigrate(db)

    if err := seedUsers(db, logger); err != nil {
        logger.Error("seed users failed", "err", err)
    }
    if err := seedPlaces(db, logger); err != nil {
        logger.Error("seed places failed", "err", err)
    }
    if err := seedBookings(db, logger); err != nil {
        logger.Error("seed bookings failed", "err", err)
    }
    if err := seedReviews(db, logger); err != nil {
        logger.Error("seed reviews failed", "err", err)
    }

    logger.Info("seeding finished")
}

func autoMigrate(db *gorm.DB) {
    _ = db.AutoMigrate(&models.User{}, &models.Place{}, &models.Booking{}, &models.Review{})
}

func seedUsers(db *gorm.DB, logger *slog.Logger) error {
    var cnt int64
    db.Model(&models.User{}).Count(&cnt)
    if cnt > 0 {
        logger.Info("users already seeded")
        return nil
    }

    users := []models.User{
        {Phone: "+70000000001", FirstName: "Иван", LastName: "Иванов"},
        {Phone: "+70000000002", FirstName: "Мария", LastName: "Петрова"},
    }

    for i := range users {
        if err := db.Create(&users[i]).Error; err != nil {
            return err
        }
        logger.Info("user created", "id", users[i].ID, "phone", users[i].Phone)
    }
    return nil
}

func seedPlaces(db *gorm.DB, logger *slog.Logger) error {
    var cnt int64
    db.Model(&models.Place{}).Count(&cnt)
    if cnt > 0 {
        logger.Info("places already seeded")
        return nil
    }

    places := []models.Place{
        {
            Name:         "Open Workspace",
            Type:         models.PlaceWorkspace,
            Description:  "Большой открытый коворкинг",
            PricePerHour: 100,
            IsActive:     true,
            CreatedAt:    time.Now(),
        },
        {
            Name:         "Meeting Room A",
            Type:         models.PlaceMeetingRoom,
            Description:  "Переговорная на 6 человек",
            PricePerHour: 500,
            IsActive:     true,
            CreatedAt:    time.Now(),
        },
    }

    for i := range places {
        if err := db.Create(&places[i]).Error; err != nil {
            return err
        }
        logger.Info("place created", "id", places[i].ID, "name", places[i].Name)
    }
    return nil
}

func seedBookings(db *gorm.DB, logger *slog.Logger) error {
    var cnt int64
    db.Model(&models.Booking{}).Count(&cnt)
    if cnt > 0 {
        logger.Info("bookings already seeded")
        return nil
    }

    var user models.User
    if err := db.First(&user).Error; err != nil {
        return nil // нет юзера — ничего не делаем
    }
    var place models.Place
    if err := db.First(&place).Error; err != nil {
        return nil
    }

    start := time.Now().Add(24 * time.Hour).Truncate(time.Minute)
    end := start.Add(2 * time.Hour)
    durationHours := end.Sub(start).Hours()
    price := durationHours * place.PricePerHour
    total := math.Round(price*100) / 100

    booking := models.Booking{
        UserID:     user.ID,
        PlaceID:    place.ID,
        StartTime:  start,
        EndTime:    end,
        TotalPrice: total,
        Status:     models.BookingActive,
    }

    if err := db.Create(&booking).Error; err != nil {
        return err
    }
    logger.Info("booking created", "id", booking.ID)
    return nil
}

func seedReviews(db *gorm.DB, logger *slog.Logger) error {
    var cnt int64
    db.Model(&models.Review{}).Count(&cnt)
    if cnt > 0 {
        logger.Info("reviews already seeded")
        return nil
    }

    var user models.User
    if err := db.First(&user).Error; err != nil {
        return nil
    }
    var place models.Place
    if err := db.First(&place).Error; err != nil {
        return nil
    }

    review := models.Review{
        UserID:     user.ID,
        PlaceID:    place.ID,
        Rating:     5,
        Text:       "Отличное место, всё понравилось.",
        IsApproved: true,
        CreatedAt:  time.Now(),
    }

    if err := db.Create(&review).Error; err != nil {
        return err
    }
    logger.Info("review created", "id", review.ID)
    return nil
}