package img

import (
	"barcelona-watch/utils"
	"image"
	"image/png"
	"net/http"
	"os"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

func downloadImage(url, filepath string) error {
	// Get the image from the URL
	response, err := http.Get(url)
	utils.HandleErr("Error:", err)
	defer response.Body.Close()

	// Create the file
	file, err := os.Create(filepath)
	utils.HandleErr("Error:", err)
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(response.Body)
	utils.HandleErr("Error:", err)

	// Save the decoded image as PNG
	err = png.Encode(file, img) // Change jpeg.Encode to png.Encode
	utils.HandleErr("Error:", err)

	return nil
}

// Function to resize an image to the specified width and height
func resizeImage(img image.Image, width, height uint) image.Image {
	// Use the resize package to resize the image while maintaining aspect ratio
	return resize.Resize(width, height, img, resize.Lanczos3)
}

// Function to create the banner with specified image paths
func createMatchBanner(homeLogoPath, awayLogoPath string) error {
	const imgWidth = 1200
	const imgHeight = 500

	// Create a new image context
	dc := gg.NewContext(imgWidth, imgHeight)

	// Set the background color to black
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Load the home team logo
	homeImg, err := gg.LoadImage(homeLogoPath)
	utils.HandleErr("Error:", err)

	// Load the away team logo
	awayImg, err := gg.LoadImage(awayLogoPath)
	utils.HandleErr("Error:", err)

	// Resize logos for better fit
	homeLogoScaledWidth := uint(300)
	resizedHomeImg := resizeImage(homeImg, homeLogoScaledWidth, 0)

	awayLogoScaledWidth := uint(300)
	resizedAwayImg := resizeImage(awayImg, awayLogoScaledWidth, 0)

	// Draw home team logo on the left (centered vertically)
	dc.DrawImageAnchored(resizedHomeImg, 300, imgHeight/2, 0.5, 0.5)

	// Draw away team logo on the right (centered vertically)
	dc.DrawImageAnchored(resizedAwayImg, imgWidth-300, imgHeight/2, 0.5, 0.5)

	// Save the output as PNG
	err = dc.SavePNG("match_banner.png")
	utils.HandleErr("Error:", err)

	return nil
}

// New function for generating the banner with provided image URLs
func GenerateBannerFromURLs(homeLogoURL, awayLogoURL string) {
	// Download the home team logo
	err := downloadImage(homeLogoURL, "home_logo.jpg")
	utils.HandleErr("Error downloading home logo:", err)

	// Download the away team logo
	err = downloadImage(awayLogoURL, "away_logo.jpg")
	utils.HandleErr("Error downloading away logo:", err)

	// Create the banner image
	err = createMatchBanner("home_logo.jpg", "away_logo.jpg")
	utils.HandleErr("Error creating banner:", err)

	// Delete the home team logo after banner creation
	err = utils.DeleteFile("home_logo.jpg")
	utils.HandleErr("Error deleting home logo:", err)

	// Delete the away team logo after banner creation
	err = utils.DeleteFile("away_logo.jpg")
	utils.HandleErr("Error deleting away logo:", err)
}
