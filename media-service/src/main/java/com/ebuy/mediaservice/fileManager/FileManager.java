package com.ebuy.mediaservice.fileManager;

import io.minio.BucketExistsArgs;
import io.minio.GetObjectArgs;
import io.minio.GetObjectResponse;
import io.minio.MakeBucketArgs;
import io.minio.MinioClient;
import io.minio.ObjectWriteResponse;
import io.minio.PutObjectArgs;
import io.minio.errors.MinioException;

import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.security.InvalidKeyException;
import java.security.NoSuchAlgorithmException;
import java.util.UUID;

import javax.imageio.stream.ImageOutputStream;

import org.checkerframework.checker.units.qual.s;
import org.springframework.core.io.InputStreamSource;
import org.springframework.web.multipart.MultipartFile;

import com.ebuy.mediaservice.imageCompressor.CompressedImageResult;
import com.ebuy.mediaservice.imageCompressor.ImageCompressor;

public class FileManager {

    private final String bucketName;
    private final MinioClient minioClient;
    private final ImageCompressor imageCompressor;

    public FileManager(String bucketName) throws IOException, NoSuchAlgorithmException, InvalidKeyException {
        MinioClient minioClient = MinioClient.builder()
                .endpoint("http://minio:9000")
                .credentials("minio", "minio123")
                .build();
        // Make 'asiatrip' bucket if not exist.
        try {
            boolean found = minioClient.bucketExists(BucketExistsArgs.builder().bucket(bucketName).build());
            if (!found) {
                // Make a new bucket called 'asiatrip'.
                minioClient.makeBucket(MakeBucketArgs.builder().bucket(bucketName).build());
            } else {
                System.out.println("Bucket 'media-service' already exists.");
            }

        } catch (MinioException e) {
            System.out.println("Error occurred: " + e);
            System.out.println("HTTP trace: " + e.httpTrace());
        }
        this.bucketName = bucketName;
        this.minioClient = minioClient;
        this.imageCompressor = new ImageCompressor();
    }

    public void uploadFile(UUID objectId, ByteArrayInputStream file, Long size, String contentType)
            throws IOException, NoSuchAlgorithmException, InvalidKeyException {
        try {
            minioClient.putObject(
                    PutObjectArgs.builder().bucket(bucketName).object(objectId.toString())
                            .stream(file,
                                    size, -1)
                            .contentType(contentType).build());
        } catch (MinioException e) {
            System.out.println("Error occurred: " + e);
            System.out.println("HTTP trace: " + e.httpTrace());
        }
    }

    public GetObjectResponse getFile(String objectId)
            throws IOException, NoSuchAlgorithmException, InvalidKeyException, MinioException {
        GetObjectResponse resp;
        try {
            resp = minioClient
                    .getObject(GetObjectArgs.builder().bucket(bucketName).object(objectId).build());
        } catch (MinioException e) {
            throw e;
        }
        return resp;
    }
}