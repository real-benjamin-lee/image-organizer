# Image Organizer

Extract images from sub-folders into a single directory

## Intro

For object detection tasks, we first need to label the images in a dataset. But in real world, the dataset provided by our clients could be nasty. Imagine your client have just sent you a 1TB folder of files which looks like this:

    report/
        2019-1-1/
            objectA.jpg
            objectB.png
            objectC.BMP
            summary.docx
        2019-1-2/
            objectA.JPEG
            further-inspection/
                objectA-1.JPEG
                objectA-2.JPEG
                objectA-3.JPEG
                summary.docx
            objectB.JPEG
            objectC.JPEG
            summary.docx
        ...
    


To solve this problem, I wrote a tiny script in [Go](https://golang.org/) to extract images from dataset sub-folders into a single directory. if you run the script by executing

    $imo -i report -o result

You would get a file structure under `./result` like this

    1.jpg // objectA.jpg
    2.png // objectB.png
    3.bmp // objectC.bmp
    4.jpeg // objectA.JPEG
    5.jpeg // objectA-1.JPEG
    6.jpeg // objectA-2.JPEG
    7.jpeg // objectA-3.JPEG
    8.jpeg // objectB.JPEG
    9.jpeg // objectC.JPEG
    ...

## Build

`clone` or `go get` this repository, `cd` to it and run `go build imo.go`

## Usage

    # organize current directory and copy images to ./image-organizer
    imo 

    # specify input & output directories
    imo -i <inputDir> -o <outputDir>

    # specify file extensions to search
    # note: file extensions would be auto-converted to lowercase
    #       which means 'jpg' would match both 'jpg' and 'JPG' 
    imo -e jpg|jpeg|bmp|png|tga

    # set search depth to 5
    imo -d 5

    # log error messages
    imo -v

    # log all messages
    imo -vv

    # show help generated by golang/pkg/flag
    imo -h
    

## License

[MIT](LICENSE.txt)

---
*Benjamin Lee, 2019*