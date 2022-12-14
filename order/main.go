package main

import (
    "github.com/gofiber/fiber/v2"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

type Order struct{
    Id          uint `json:"id"`
    Product_id  uint `json:"product_id"`   
    Qty         uint `json:"qty"`
    Price       uint `json:"price"`
    Payment     uint `json:"payment"`
    Status      uint `json:"status`
}

type Product struct{
    Id            uint      `json:"id"`
    Name_product  string    `json:"name_product"`
    Amount        uint      `json:"amount"`
    Price         uint      `json:"price"`
}

type Cart struct{
    Id          uint `json:"id"`
    Product     Product `json:"product"`
    Qty         uint `json:"qty"`
    Price       uint `json:"price"`
}

type ResponseOk struct{
    Status string `json:"status"`
    Message string `json:"message"`
}

type ResponseData struct{
    Status string `json:"status"`
    Data []Product `json:"data"`
}


func main() {
    db, err := gorm.Open(mysql.Open("root:farhan123@tcp(127.0.0.1:3306)/order_ms"), &gorm.Config{})

    if err != nil{
        panic(err)
    }

    db.AutoMigrate(Order{})

    app := fiber.New()

    app.Use(cors.New())

    app.Get("/api/cart", func(c *fiber.Ctx) error {
        var cart []Cart

        db.Find(&cart)

        c.Set("Content-type", "application/json")

        resp := ResponseData{Status:"200", Data:products}

        return c.JSON(resp)
    })

    app.Get("/api/orders", func(c *fiber.Ctx) error {
        var order []Order

        db.Find(&order)

        c.Set("Content-type", "application/json")

        resp := ResponseData{Status:"200", Data:products}

        return c.JSON(resp)
    })

    app.Listen(":3001")
}