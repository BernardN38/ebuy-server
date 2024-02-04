package com.ebuy.mediaservice.imageCompressor;

import java.util.Iterator;

import javax.imageio.IIOImage;
import javax.imageio.ImageIO;
import javax.imageio.ImageWriteParam;
import javax.imageio.ImageWriter;
import javax.imageio.stream.ImageOutputStream;

import java.awt.image.BufferedImage;
import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.ByteArrayOutputStream;

public class ImageCompressor {
    private static final ImageWriter imageWriter = ImageIO.getImageWritersByFormatName("jpg").next();

    public ImageCompressor() {

    }

    public static CompressedImageResult compressImage(ByteArrayInputStream image) throws IOException {
        // Record the start time
        long startTime = System.currentTimeMillis();
        try (ByteArrayOutputStream bos = new ByteArrayOutputStream()) {
            // Read the image using ImageIO
            BufferedImage bufferedImage = ImageIO.read(image);

            try (ImageOutputStream outputStream = ImageIO.createImageOutputStream(bos)) {
                imageWriter.setOutput(outputStream);

                // Set compression parameters
                ImageWriteParam params = imageWriter.getDefaultWriteParam();
                params.setCompressionMode(ImageWriteParam.MODE_EXPLICIT);
                params.setCompressionQuality(0.5f);

                // Write the compressed image
                imageWriter.write(null, new IIOImage(bufferedImage, null, null), params);
            }
            // Record the end time
            long endTime = System.currentTimeMillis();

            // Calculate and print the execution time
            long executionTime = endTime - startTime;
            System.out.println("Execution time: " + executionTime + " milliseconds");
            return new CompressedImageResult(bos, bos.size());
        }

    }
}
