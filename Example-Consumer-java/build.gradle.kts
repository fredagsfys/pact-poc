import io.ktor.http.*

plugins {
    java
    id("au.com.dius.pact") version "4.3.10"
    id("de.undercouch.download") version "5.6.0"
    id("org.openapi.generator") version "6.6.0"
}

group = "org.example.consumer"
version = "1.0-SNAPSHOT"
val osaPlatformVersion = "2.0.0"

java {
    toolchain {
        languageVersion = JavaLanguageVersion.of(21)
    }
}

dependencies {
    implementation("com.squareup.retrofit2:retrofit:2.9.0")
    implementation("com.squareup.retrofit2:converter-jackson:2.9.0")
    implementation("org.slf4j:slf4j-api:2.0.16")

    // test
    testImplementation("org.json:json:20220320")
    testImplementation("org.junit.jupiter:junit-jupiter-api:5.8.1")
    testRuntimeOnly("org.junit.jupiter:junit-jupiter-engine:5.8.1")
    testImplementation("au.com.dius.pact.consumer:junit5:4.3.10")
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

tasks.register("downloadOSA") {
    doLast {
        val targetUrl = "https://api.github.com/repos/OWNER/REPO/contents/specifications/resources/.../$osaPlatformVersion/ping_$osaPlatformVersion.yaml"
        val destinationFile = layout.buildDirectory.file("${layout.buildDirectory}/tmp/ping_$osaPlatformVersion.yaml").get().asFile
        download.run {
            src(targetUrl)
            dest(destinationFile)

            val ghToken = System.getenv("GITHUB_TOKEN")
            if (ghToken.isNotEmpty()) {
                header("Accept", "application/vnd.github.raw+json")
                header("Authorization", "Bearer $ghToken")
                header("X-GitHub-Api-Version", "2022-11-28")
            }
        }
    }
}
defaultTasks("downloadOSA")

openApiGenerate {
    generatorName = "spring"
    inputSpec = "${layout.buildDirectory.get()}/tmp/ping_$osaPlatformVersion.yaml"
    outputDir = "${layout.buildDirectory.get()}/generated"
    apiPackage = "com.platform.generated.api"
    modelPackage = "com.platform.generated.model"
}

sourceSets {
    main {
        java {
            srcDir("$${layout.buildDirectory.get()}/generated/src/main/java")
        }
    }
}