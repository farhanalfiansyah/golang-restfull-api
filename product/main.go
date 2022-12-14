package main

import (
    "github.com/go-playground/validator/v10"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

type Product struct{
    Id            uint      `json:"id"`             
    Name_product  string    `json:"name_product"    validate:"required"` 
    Amount        int      `json:"amount"           validate:"required,number"`         
    Price         int      `json:"price"            validate:"required,number"`          
}   

type ResponseOk struct{
    Status string `json:"status"`
    Message string `json:message`
}

type ResponseData struct{
    Status string `json:"status"`
    Data []Product `json:"data"`
}

var validate = validator.New()

func main() {
    db, err := gorm.Open(mysql.Open("root:farhan123@tcp(127.0.0.1:3306)/product_ms"), &gorm.Config{})

    if err != nil{
        panic(err)
    }

    db.AutoMigrate(Product{})

    app := fiber.New()

    app.Use(cors.New())

    app.Get("/api/product", func(c *fiber.Ctx) error {
        var products []Product

        db.Find(&products)

        c.Set("Content-type", "application/json")

        resp := ResponseData{Status:"success", Data:products}

        return c.JSON(resp)
    })

    app.Get("/api/product/:id", func(c *fiber.Ctx) error {
        var product Product
        id := c.Params("id")
        result := db.Find(&product, id)
        if result.RowsAffected == 0 {
            return c.Status(400).JSON(fiber.Map{
                "status" : "failed",
                "message" : "Product id not found!",
            })
        } 
        return c.Status(200).JSON(&product)
    })

    app.Post("/api/product", func(c *fiber.Ctx) error {
        var product Product
        if err := c.BodyParser(&product); err != nil{
            return c.Status(400).JSON(fiber.Map{
                "status" : "failed",
                "message" : err.Error(),
            })
        }

        errValidate := validate.Struct(product)
        if errValidate != nil{
            return c.Status(400).JSON(fiber.Map{
                "status" : "failed",
                "message" : errValidate.Error(),
            })
        }

        db.Create(&product)
        c.Set("Content-Type", "application/json")

        resp := ResponseOk{Status:"Success", Message:"Successfully Created Product!"}

        return c.JSON(resp)
    })

    app.Delete("/api/product/:id", func(c *fiber.Ctx) error {
        var product Product

        result := db.Delete(&product, c.Params("id"))
        c.Set("Content-Type", "application/json")
        if result.RowsAffected == 0 {
            return c.Status(400).JSON(fiber.Map{
                "status" : "failed",
                "message" : "Product id not found!",
            })
        }
        resp := ResponseOk{Status:"Success", Message:"Successfully Delete Product!"}

        return c.JSON(resp)
    })

    app.Put("/api/product/:id", func(c *fiber.Ctx) error{
        var product Product
        id := c.Params("id")

        if err := c.BodyParser(&product); err != nil{
            return c.Status(400).JSON(fiber.Map{
                "status" : "failed",
                "message" : err.Error(),
            })
        }
        db.Where("id = ?", id).Updates(&product)

        c.Set("Content-Type", "application/json")

        resp := ResponseOk{Status:"Success", Message:"Successfully Update Product!"}

        return c.JSON(resp)
    })

    app.Listen(":3000")
}