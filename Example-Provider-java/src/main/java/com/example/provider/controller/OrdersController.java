package com.example.provider.controller;

import com.example.provider.model.Order;
import com.example.provider.repository.OrdersRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import java.io.IOException;
import java.util.List;

@RestController
public class OrdersController {
    private final OrdersRepository ordersRepository;

    @Autowired
    public OrdersController(OrdersRepository ordersRepository) {
        this.ordersRepository = ordersRepository;
    }

    @GetMapping("/orders")
    List<Order> getAllOrders() throws IOException {
        return ordersRepository.getOrders();
    }
}
