package com.ebuy.mediaservice.messaging;

import java.io.IOException;
import java.util.Optional;

import org.springframework.amqp.core.Message;
import org.springframework.amqp.core.MessageListener;
import org.springframework.amqp.support.converter.MessageConversionException;
import org.springframework.amqp.support.converter.MessageConverter;

import com.ebuy.mediaservice.UserMediaRepository.UserMediaRepository;
import com.ebuy.mediaservice.entities.UserMedia.UserMedia;
import com.fasterxml.jackson.databind.ObjectMapper;

public class MessageReceiver implements MessageListener {

    private final MessageConverter messageConverter;// Create ObjectMapper
    private final ObjectMapper objectMapper = new ObjectMapper();
    private final UserMediaRepository repository;

    public MessageReceiver(MessageConverter messageConverter, UserMediaRepository userMediaRepository) {
        this.messageConverter = messageConverter;
        this.repository = userMediaRepository;
    }

    public void receive(byte[] message) {
        System.out.println(new String(message));
    }

    public void onMessage() {

    }

    private void handleImageUploadedDataMessage(ImageCompressedMessage message) {
        System.out.println("Received ImageUploadedData: " + message);
        // Add your processing logic for ImageUploadedData
    }

    @Override
    public void onMessage(Message message) {
        try {
            ImageCompressedMessage imageCompressedMessage = objectMapper.readValue(message.getBody(),
                    ImageCompressedMessage.class);
            Optional<UserMedia> userMedia = repository.findById(imageCompressedMessage.getMediaId());
            if (!userMedia.isPresent()) {
                return;
            }
            userMedia.get().setCompression_status(true);
            repository.save(userMedia.get());
            System.out
                    .println(
                            "Media compression status updated, new status: " + userMedia.get().getCompression_status()
                                    + imageCompressedMessage.getUserId());
        } catch (MessageConversionException | IOException e) {
            System.err.println("Error converting message: " + e.getMessage());
        }
        System.out.println("message" + message);
    }
}
