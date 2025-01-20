plugins {
	java
	id("org.springframework.boot") version "3.4.1"
	id("io.spring.dependency-management") version "1.1.7"
}

group = "com.example.provider"
version = "0.0.1-SNAPSHOT"

java {
	toolchain {
		languageVersion = JavaLanguageVersion.of(21)
	}
}

repositories {
	mavenCentral()
}

dependencies {
	implementation("org.springframework.boot:spring-boot-starter")
	implementation("org.springframework.boot:spring-boot-starter-web")

	// test
	testImplementation("org.springframework.boot:spring-boot-starter-test")
	testRuntimeOnly("org.junit.platform:junit-platform-launcher")
	testImplementation("commons-codec:commons-codec")
	testImplementation("au.com.dius.pact.provider:junit5spring:4.3.7")
}

tasks.withType<Test> {
	useJUnitPlatform()

	systemProperty("pact.provider.version", project.version)
	systemProperty("pact.verifier.publishResults", "true")
	systemProperty("pact_do_not_track", "true")
}

