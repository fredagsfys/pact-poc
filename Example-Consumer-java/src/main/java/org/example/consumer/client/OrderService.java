package org.example.consumer.client;

import org.example.consumer.model.Order;
import retrofit2.Call;
import retrofit2.http.GET;
import java.util.List;

public interface OrderService {
    @GET("orders")
    Call<List<Order>> getOrders();
}