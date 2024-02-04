FROM eclipse-temurin:21-jdk-alpine

RUN apk --no-cache add curl
VOLUME /tmp
ARG JAR_FILE
COPY ./target/media-service-0.0.1-SNAPSHOT.jar app.jar

ENTRYPOINT ["java", "-Xmx1g", "-jar","/app.jar"]