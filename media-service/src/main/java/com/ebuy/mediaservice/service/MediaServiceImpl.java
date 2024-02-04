package com.ebuy.mediaservice.service;

import org.bouncycastle.jcajce.provider.asymmetric.ec.SignatureSpi.ecCVCDSA3_224;
import org.hibernate.ObjectNotFoundException;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.core.io.InputStreamResource;
import org.springframework.stereotype.Service;
import org.springframework.web.multipart.MultipartFile;

import com.ebuy.mediaservice.UserMediaRepository.UserMediaRepository;
import com.ebuy.mediaservice.entities.UserMedia.UserMedia;
import com.ebuy.mediaservice.fileManager.FileManager;
import com.ebuy.mediaservice.messaging.ImageUploadedMessage;

import io.minio.GetObjectResponse;

import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.security.InvalidKeyException;
import java.security.NoSuchAlgorithmException;
import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

import javax.naming.NameNotFoundException;

@Service
public class MediaServiceImpl implements MediaService {
    private final FileManager fileManager;
    private final UserMediaRepository repository;
    private final RabbitTemplate rabbitTemplate;

    public MediaServiceImpl(FileManager fileUploaderBean, UserMediaRepository userMediaRepositoryBean,
            RabbitTemplate rabbitTemplate)
            throws InvalidKeyException, NoSuchAlgorithmException, IOException {
        this.fileManager = fileUploaderBean;
        this.repository = userMediaRepositoryBean;
        this.rabbitTemplate = rabbitTemplate;
    }

    @Override
    public Optional<GetObjectResponse> GetImage(String imageId) throws Exception {
        GetObjectResponse obj;
        try {
            obj = fileManager.getFile(imageId);
        } catch (Exception e) {
            throw e;
        }
        return Optional.of(obj);
    }

    @Override
    public Long CreateUserMedia(MultipartFile image, Long userId) throws Exception {
        UUID genertatedIdFull = UUID.randomUUID();
        UUID genertatedIdCompressed = UUID.randomUUID();
        try {
            fileManager.uploadFile(genertatedIdFull, new ByteArrayInputStream(image.getBytes()), image.getSize(),
                    image.getContentType());
            UserMedia userMedia = new UserMedia(userId, genertatedIdFull, genertatedIdCompressed);
            repository.save(userMedia);
            ImageUploadedMessage message = new ImageUploadedMessage(userMedia.getId(), genertatedIdFull,
                    genertatedIdCompressed,
                    image.getContentType());
            rabbitTemplate.convertAndSend("media_events", "upload", message);
            return userMedia.getId();
        } catch (InvalidKeyException e) {
            // TODO Auto-generated catch block
            e.printStackTrace();
        } catch (NoSuchAlgorithmException e) {
            // TODO Auto-generated catch block
            e.printStackTrace();
        } catch (IOException e) {
            // TODO Auto-generated catch block
            e.printStackTrace();
        }
        return Long.valueOf(0);
    }

    @Override
    public List<UserMedia> GetAll() throws Exception {
        Iterable<UserMedia> all = repository.findAll();

        // Using loop to convert Iterable to List
        List<UserMedia> userList = new ArrayList<>();
        for (UserMedia user : all) {
            userList.add(user);
        }

        return userList;
    }

}
