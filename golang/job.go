package api

import (
	"image"
	"image/color"
	"io"
	"math"
	"mime/multipart"
	"os"

	"github.com/disintegration/imaging"
)

func putInSquare(img image.Image, size int) *image.NRGBA {
	dst := imaging.New(size, size, color.RGBA{255, 255, 255, 0})
	return imaging.Paste(dst, img, image.Pt(0, 0))
}

func performJobAllThree(dbw *DBWorker, job *Job, mediaRoot string, fileHeader *multipart.FileHeader) error {
	if err := performJobOrig(dbw, job, mediaRoot, fileHeader); err != nil {
		return err
	}

	if err := performJobSquareOrig(dbw, job, mediaRoot, fileHeader); err != nil {
		return err
	}

	if err := performJobSquareSmall(dbw, job, mediaRoot, fileHeader); err != nil {
		return err
	}
	return nil
}

func performJobSquareSmall(dbw *DBWorker, job *Job, mediaRoot string, fileHeader *multipart.FileHeader) error {
	dbImg, err := NewImageFromForm(job, mediaRoot, fileHeader)
	if err != nil {
		return err
	}

	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()
	src, err := imaging.Decode(file)
	if err != nil {
		return err
	}
	res := imaging.Crop(src, image.Rect(0, 0, 256, 256))
	b := res.Bounds()
	if b.Max.X < 256 || b.Max.Y < 256 {
		res = putInSquare(res, 256)
	}
	if err := imaging.Save(res, dbImg.Path); err != nil {
		return err
	}
	stat, err := os.Stat(dbImg.Path)
	if err != nil {
		return err
	}
	dbImg.Width = 256
	dbImg.Height = 256
	dbImg.Size = int64(stat.Size())
	return dbw.SaveNewImage(dbImg)
}

func performJobSquareOrig(dbw *DBWorker, job *Job, mediaRoot string, fileHeader *multipart.FileHeader) error {
	dbImg, err := NewImageFromForm(job, mediaRoot, fileHeader)
	if err != nil {
		return err
	}

	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()
	src, err := imaging.Decode(file)
	if err != nil {
		return err
	}
	b := src.Bounds()
	size := math.Max(float64(b.Max.X), float64(b.Max.Y))
	dbImg.Width = int(size)
	dbImg.Height = int(size)
	img := putInSquare(src, int(size))
	if err := imaging.Save(img, dbImg.Path); err != nil {
		return err
	}
	stat, err := os.Stat(dbImg.Path)
	if err != nil {
		return err
	}
	dbImg.Size = int64(stat.Size())
	return dbw.SaveNewImage(dbImg)
}

func performJobOrig(dbw *DBWorker, job *Job, mediaRoot string, fileHeader *multipart.FileHeader) error {
	dbImg, err := NewImageFromForm(job, mediaRoot, fileHeader)
	if err != nil {
		return err
	}

	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()
	src, err := imaging.Decode(file)
	if err != nil {
		return err
	}
	b := src.Bounds()
	dbImg.Width = b.Max.X
	dbImg.Height = b.Max.Y

	out, err := os.Create(dbImg.Path)
	if err != nil {
		return err
	}
	defer out.Close()
	file.Seek(0, io.SeekStart)
	_, err = io.Copy(out, file)
	if err != nil {
		return err
	}
	return dbw.SaveNewImage(dbImg)
}
