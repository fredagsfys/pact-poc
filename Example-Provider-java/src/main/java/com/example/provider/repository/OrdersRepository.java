package com.example.provider.repository;

import com.example.provider.model.Order;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.stereotype.Component;

import java.io.IOException;
import java.net.URL;
import java.util.List;

@Component
public class OrdersRepository {
    private final ObjectMapper objectMapper = new ObjectMapper();

    public List<Order> getOrders() throws IOException {
        URL resource = getClass().getClassLoader().getResource("orders.json");
        return objectMapper.readValue(resource, new TypeReference<>() {
        });
    }
}
