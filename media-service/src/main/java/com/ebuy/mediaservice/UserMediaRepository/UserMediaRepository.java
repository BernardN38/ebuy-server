package com.ebuy.mediaservice.UserMediaRepository;

import java.util.Optional;
import java.util.UUID;

import org.springframework.data.repository.CrudRepository;

import com.ebuy.mediaservice.entities.UserMedia.UserMedia;

public interface UserMediaRepository extends CrudRepository<UserMedia, Long> {
    Optional<UserMedia> findById(Long id);

    Optional<UserMedia> findByMediaId(UUID media_id);
}