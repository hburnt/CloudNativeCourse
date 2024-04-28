package recipeinfoapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
//	"net/url"
	"time"
)

type Ingredient struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
}

type RecipeID int64

type Recipe struct {
	ID          int          `json:"id"`
	Title       string       `json:"title"`
	Image       string       `json:"image"`
	Ingredients []Ingredient `json:"extendedIngredients"`
}

type RecipeInfo struct {
	ID          int
	Title       string
	Ingredients []Ingredient
}

type SpoonRecipeResponse struct {
	ID          int          `json:"id"`
	Title       string       `json:"title"`
	Ingredients []Ingredient `json:"extendedIngredients"`
}

type Client struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(key string) *Client {
	return &Client{
		APIKey: key,
		BaseURL: "https://api.spoonacular.com",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) FormatRecipeInfoURL(recipeID int) string {
	return fmt.Sprintf("%s/recipes/%d/information?includeNutrition=false&apiKey=%s", c.BaseURL, recipeID, c.APIKey)
}

func (c *Client) GetRecipe(recipeID int) (RecipeInfo, error) {
	URL := c.FormatRecipeInfoURL(recipeID)
	resp, err := c.HTTPClient.Get(URL)
	if err != nil {
		return RecipeInfo{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return RecipeInfo{}, fmt.Errorf("could not find recipe with ID: %d", recipeID)
	}
	if resp.StatusCode != http.StatusOK {
		return RecipeInfo{}, fmt.Errorf("unexpected response status %q", resp.Status)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return RecipeInfo{}, err
	}
	recipe_info, err := parseRecipeInfoResponse(data)
	if err != nil {
		return RecipeInfo{}, err
	}
	return recipe_info, nil
}

func parseRecipeInfoResponse(data []byte) (RecipeInfo, error) {
	var resp SpoonRecipeResponse
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return RecipeInfo{}, fmt.Errorf("invalid API response %s: %w", data, err)
	}

	ingredients := make([]Ingredient, len(resp.Ingredients))
	for i, ing := range resp.Ingredients {
		ingredient := Ingredient{
			ID:     ing.ID,
			Name:   ing.Name,
			Amount: ing.Amount,
		}
		ingredients[i] = ingredient
	}

	info := RecipeInfo{
		ID:          resp.ID,
		Title:       resp.Title,
		Ingredients: ingredients,
	}
	return info, nil
}
