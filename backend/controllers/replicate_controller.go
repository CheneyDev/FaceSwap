package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type ReplicateInput struct {
	Image         string `json:"image"`
	ImageToBecome string `json:"image_to_become"`
}

type ReplicateRequest struct {
	Version string         `json:"version"`
	Input   ReplicateInput `json:"input"`
}

type ReplicateResponse struct {
	ID          string         `json:"id"`
	Model       string         `json:"model"`
	Version     string         `json:"version"`
	Input       ReplicateInput `json:"input"`
	Logs        string         `json:"logs"`
	Output      []string       `json:"output"`
	Error       string         `json:"error"`
	Status      string         `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	StartedAt   time.Time      `json:"started_at"`
	CompletedAt time.Time      `json:"completed_at"`
	URLs        struct {
		Cancel string `json:"cancel"`
		Get    string `json:"get"`
	} `json:"urls"`
	Metrics struct {
		PredictTime float64 `json:"predict_time"`
	} `json:"metrics"`
}

func CreateReplicatePrediction(c *gin.Context) {
	var input ReplicateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reqBody, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "https://api.replicate.com/v1/predictions", bytes.NewBuffer(reqBody))
	replicateToken := "Bearer " + os.Getenv("REPLICATE_API_TOKEN")
	req.Header.Set("Authorization", replicateToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	var prediction ReplicateResponse
	if err := json.NewDecoder(resp.Body).Decode(&prediction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, prediction)
}

func GetReplicatePrediction(c *gin.Context) {
	id := c.Param("id")
	url := fmt.Sprintf("https://api.replicate.com/v1/predictions/%s", id)
	replicateToken := "Bearer " + os.Getenv("REPLICATE_API_TOKEN")
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", replicateToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var prediction ReplicateResponse
	if err := json.Unmarshal(body, &prediction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, prediction)
}
