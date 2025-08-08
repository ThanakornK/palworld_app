package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"palworld_tools/dto"
	"palworld_tools/models"
	"palworld_tools/services/datamanage"
	"palworld_tools/services/options"
	"palworld_tools/services/scrapper"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func main() {
	// var functionName string

	// menus := map[string]string{
	// 	"1": "Update Data",
	// 	"2": "Add Pal to store",
	// }

	// fmt.Println("Enter function: ")
	// for i, menu := range menus {
	// 	fmt.Printf("%s. %s\n", i, menu)
	// }
	// // fmt.Println("1. Update Data")
	// // fmt.Println("2. Add Pal to store")
	// fmt.Print("Function number: ")
	// reader := bufio.NewReader(os.Stdin)
	// functionName, _ = reader.ReadString('\n')
	// functionName = strings.TrimSpace(functionName)

	// switch functionName {
	// case "1":
	// 	fmt.Println("Update Data")
	// 	loading(updateData)

	// case "2":
	// 	// loop enter until enter empty
	// 	for {
	// 		fmt.Println("Add Pal to store")
	// 		err := AddPalToStore()
	// 		if err != nil {
	// 			fmt.Println("Error Add Pal to store: ", err)
	// 		}

	// 		fmt.Println("Add another pal? (y/n)")
	// 		reader := bufio.NewReader(os.Stdin)
	// 		answer, _ := reader.ReadString('\n')
	// 		answer = strings.TrimSpace(answer)
	// 		if answer == "n" {
	// 			break
	// 		}

	// 	}

	// default:
	// 	fmt.Println("Invalid function number")
	// }

	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:3001"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	r.GET("/update-data", func(ctx *gin.Context) {
		err := updateData()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Data updated successfully"})

	})

	r.POST("/add-pal", func(ctx *gin.Context) {
		var pal dto.AddPalRequest

		if err := ctx.ShouldBindJSON(&pal); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := datamanage.AddPal(pal.Name, pal.Gender, pal.PassiveSkills)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Pal added successfully"})
	})

	r.GET("/store", func(ctx *gin.Context) {
		result, err := datamanage.ReadStoredPals()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return

		}

		palDex, err := datamanage.ReadPaldex()

		// make map of palDex
		palDexMap := make(map[string]models.Pal)
		for _, pal := range palDex {
			spew.Dump("CHECK PALDEX", pal.ImageUrl)

			palDexMap[strings.ToLower(pal.Name)] = pal
		}

		var pals []dto.Pal
		for _, species := range result {
			for _, pal := range species.StoredPals {

				var passiveSkills []dto.PassiveSkill
				for _, skill := range pal.PassiveSkills {
					passiveSkills = append(passiveSkills, dto.PassiveSkill{
						Name: skill,
					})
				}

				pals = append(pals, dto.Pal{
					Id:            pal.ID,
					Name:          species.Name,
					ImageUrl:      palDexMap[strings.ToLower(species.Name)].ImageUrl,
					Gender:        pal.Gender,
					PassiveSkills: passiveSkills,
				})
			}
		}

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": pals})
	})

	optionGroup := r.Group("/options")
	{
		optionGroup.GET("/passive-skills", func(ctx *gin.Context) {
			result := options.GetPassiveSkills()

			var passiveSkills []string
			passiveSkills = append(passiveSkills, result...)

			ctx.JSON(http.StatusOK, gin.H{"message": passiveSkills})
		})

		optionGroup.GET("/pal-species", func(ctx *gin.Context) {

			result := options.GetPalSpecies()

			var palSpecies []string
			palSpecies = append(palSpecies, result...)

			ctx.JSON(http.StatusOK, gin.H{"message": palSpecies})

		})
	}

	r.Run(":8080")

}

func loading(function func() error) {
	// Spinner characters
	spinner := []string{"|", "/", "-", "\\"}
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return // Stop the spinner loop when done is received
			case <-time.After(100 * time.Millisecond):
				// Print the spinner with a carriage return to overwrite the last character
				fmt.Print("\rLoading... ", spinner[time.Now().Unix()%4])
			}
		}
	}()

	err := function()
	if err != nil {
		errMessage := fmt.Sprintf("Error: %v", err)
		panic(errMessage)
	}

	done <- true // Send a signal to stop the spinner loop
}

func updateData() error {
	err := scrapper.ScrapperPalInfo()
	if err != nil {
		return err
	}

	// Wait for 2 seconds before running the next function
	time.Sleep(5 * time.Second)

	err = scrapper.ScrapperPassiveSkill()
	if err != nil {
		return err
	}

	// Wait for 2 seconds before running the next function
	time.Sleep(5 * time.Second)

	err = scrapper.BestComboPassiveSkill()
	if err != nil {
		return err
	}

	return nil
}

func AddPalToStore() error {
	addPalReader := bufio.NewReader(os.Stdin)
	fmt.Print("Pal name: ")
	palName, _ := addPalReader.ReadString('\n')
	palName = strings.TrimSpace(palName)

	fmt.Print("Pal gender (m/f): ")
	palGender, _ := addPalReader.ReadString('\n')
	palGender = strings.TrimSpace(palGender)
	if palGender != "m" && palGender != "f" {
		return fmt.Errorf("invalid pal gender")
	}

	passiveSkills := make([]string, 0)
	fmt.Println("Please enter passive skill, enter empty to done input")
	for i := 0; i < 4; i++ {
		fmt.Printf("Passive skill %d: ", i+1)
		passiveSkillName, _ := addPalReader.ReadString('\n')
		passiveSkillName = strings.TrimSpace(passiveSkillName)
		if passiveSkillName == "" {
			break
		}
		passiveSkills = append(passiveSkills, passiveSkillName)
	}

	fmt.Println("Input is done")

	err := datamanage.AddPal(palName, palGender, passiveSkills)
	if err != nil {
		return err
	}

	return nil
}

// func init() {
// 	r := gin.Default()

// 	r.GET("/update-data", func(ctx *gin.Context) {
// 		err := updateData()
// 		if err != nil {
// 			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 		ctx.JSON(http.StatusOK, gin.H{"message": "Data updated successfully"})

// 	})

// 	r.POST("/add-pal", func(ctx *gin.Context) {
// 		var pal dto.AddPalRequest

// 		if err := ctx.ShouldBindJSON(&pal); err != nil {
// 			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
// 		err := datamanage.AddPal(pal.Name, pal.Gender, pal.PassiveSkills)
// 		if err != nil {
// 			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 		ctx.JSON(http.StatusOK, gin.H{"message": "Pal added successfully"})
// 	})

// 	r.GET("/store", func(ctx *gin.Context) {
// 		result, err := datamanage.ReadStoredPals()
// 		if err != nil {
// 			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return

// 		}

// 		palDex, err := datamanage.ReadPaldex()

// 		// make map of palDex
// 		palDexMap := make(map[string]models.Pal)
// 		for _, pal := range palDex {
// 			palDexMap[pal.Name] = pal
// 		}

// 		var pals []dto.Pal
// 		for _, species := range result {
// 			for _, pal := range species.StoredPals {

// 				var passiveSkills []dto.PassiveSkill
// 				for _, skill := range pal.PassiveSkills {
// 					passiveSkills = append(passiveSkills, dto.PassiveSkill{
// 						Name: skill,
// 					})
// 				}

// 				pals = append(pals, dto.Pal{
// 					Id:            pal.ID,
// 					Name:          species.Name,
// 					ImageUrl:      palDexMap[species.Name].ImageUrl,
// 					Gender:        pal.Gender,
// 					PassiveSkills: passiveSkills,
// 				})
// 			}
// 		}

// 		if err != nil {
// 			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}

// 		ctx.JSON(http.StatusOK, gin.H{"message": pals})
// 	})

// 	optionGroup := r.Group("/options")
// 	{
// 		optionGroup.GET("/passive-skills", func(ctx *gin.Context) {
// 			result := options.GetPassiveSkills()

// 			var passiveSkills []string
// 			passiveSkills = append(passiveSkills, result...)

// 			ctx.JSON(http.StatusOK, gin.H{"message": passiveSkills})
// 		})

// 		optionGroup.GET("/pal-species", func(ctx *gin.Context) {

// 			result := options.GetPalSpecies()

// 			var palSpecies []string
// 			palSpecies = append(palSpecies, result...)

// 			ctx.JSON(http.StatusOK, gin.H{"message": palSpecies})

// 		})
// 	}

// 	r.Run(":8080")

// }
