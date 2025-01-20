package com.example.provider;

import au.com.dius.pact.provider.junit5.PactVerificationContext;
import au.com.dius.pact.provider.junitsupport.Provider;
import au.com.dius.pact.provider.junitsupport.State;
import au.com.dius.pact.provider.junitsupport.loader.PactBroker;
import au.com.dius.pact.provider.spring.junit5.PactVerificationSpringProvider;
import com.example.provider.repository.OrdersRepository;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.jupiter.api.TestTemplate;
import org.junit.jupiter.api.extension.ExtendWith;
import org.springframework.boot.test.context.SpringBootTest;
import org.mockito.Mockito;
import org.springframework.test.context.bean.override.mockito.MockitoBean;
import org.springframework.test.context.junit.jupiter.SpringExtension;
import com.example.provider.model.Order;

import java.io.IOException;
import java.net.URL;
import java.util.List;

@ExtendWith(SpringExtension.class)
@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.DEFINED_PORT)
@Provider("order_provider_java")
@PactBroker
public class ContractVerificationTest {
    @MockitoBean
    private OrdersRepository ordersRepository;
    private final ObjectMapper objectMapper = new ObjectMapper();

    @TestTemplate
    @ExtendWith(PactVerificationSpringProvider.class)
    void pactVerificationTestTemplate(PactVerificationContext context) {
        context.verifyInteraction();
    }

    @State("there are orders")
    public void thereAreOrders() throws IOException {
        Mockito.when(ordersRepository.getOrders()).thenReturn(getOrdersFromFile("orders.json"));
    }

    @State("there are no orders")
    public void thereAreNoOrders() throws IOException {
        Mockito.when(ordersRepository.getOrders()).thenReturn(getOrdersFromFile("no_orders.json"));
    }

    private List<Order> getOrdersFromFile(String filename) throws IOException {
        URL resource = getClass().getClassLoader().getResource(filename);
        return objectMapper.readValue(resource, new TypeReference<>() {
        });
    }
}
