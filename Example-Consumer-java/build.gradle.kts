import io.ktor.http.*

plugins {
    id("java")
    id( "au.com.dius.pact") version ("4.3.6")
}

group = "org.example.consumer"
version = "1.0-SNAPSHOT"

repositories {
    mavenCentral()
}

dependencies {
    implementation("com.squareup.retrofit2:retrofit:2.9.0")
    implementation("com.squareup.retrofit2:converter-jackson:2.9.0")
    implementation("org.slf4j:slf4j-api:2.0.16")

    // test
    testImplementation("org.json:json:20220320")
    testImplementation("org.junit.jupiter:junit-jupiter-api:5.8.1")
    testRuntimeOnly("org.junit.jupiter:junit-jupiter-engine:5.8.1")
    testImplementation("au.com.dius.pact.consumer:junit5:4.3.7")
    testImplementation("org.slf4j:slf4j-simple:2.0.16")
}

pact {
    publish {
        pactBrokerUrl = "http://localhost:9292/"
        tags = listOf("dev")
    }
    broker {
        pactBrokerUrl = "http://localhost:9292/"
        retryCountWhileUnknown = 3
        retryWhileUnknownInterval = 120 // 2 minutes between retries
    }
}

tasks.getByName<Test>("test") {
    useJUnitPlatform()

    systemProperty("pact.consumer.version", project.version)
    systemProperty("pact.verifier.publishResults", "true")
    systemProperty("pact_do_not_track", "true")
}