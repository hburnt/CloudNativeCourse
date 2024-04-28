package recipeapi

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "time"
)

type RecipeID int64

type Recipe struct {
    ID    int    `json:"id"`
    Title string `json:"title"`
    Image string `json:"image"`
}

type RecipeInfo struct {
    ID         int
    Title      string
    Ingredient string
}

type SpoonRecipeResponse struct {
    Results []struct {
        ID    int    `json:"id"`
        Title string `json:"title"`
    } `json:"results"`
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

func (c *Client) FormatRecipeSearchURL(search string) string {
    search = url.QueryEscape(search)
    return fmt.Sprintf("%s/recipes/complexSearch?query=%s&number=1&apiKey=%s", c.BaseURL, search, c.APIKey)
}

func (c *Client) GetRecipe(search string) (RecipeInfo, error) {
    URL := c.FormatRecipeSearchURL(search)
    resp, err := c.HTTPClient.Get(URL)
    if err != nil {
        return RecipeInfo{}, err
    }
    defer resp.Body.Close()
    if resp.StatusCode == http.StatusNotFound {
        return RecipeInfo{}, fmt.Errorf("could not find search: %s ", search)
    }
    if resp.StatusCode != http.StatusOK {
        return RecipeInfo{}, fmt.Errorf("unexpected response status %q", resp.Status)
    }
    data, err := io.ReadAll(resp.Body)
    if err != nil {
        return RecipeInfo{}, err
    }
    recipe, err := parseRecipeResponse(data)
    if err != nil {
        return RecipeInfo{}, err
    }
    return recipe, nil
}

func parseRecipeResponse(data []byte) (RecipeInfo, error) {
    var resp SpoonRecipeResponse
    err := json.Unmarshal(data, &resp)
    if err != nil {
        return RecipeInfo{}, fmt.Errorf("invalid API response %s: %w", data, err)
    }
    if len(resp.Results) == 0 {
        return RecipeInfo{}, fmt.Errorf("no recipes found")
    }
    recipe := resp.Results[0]

    info := RecipeInfo{
        ID:         recipe.ID,
        Title:      "Top Recipe",
        Ingredient: recipe.Title,
    }
    return info, nil
}
