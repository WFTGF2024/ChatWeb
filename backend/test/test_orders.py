# tests/test_orders.py
import pytest


def test_order_flow(http, base_url, test_user):
    user_id = test_user["user_id"]

    # 先创建订单
    order_payload = {
        "user_id": user_id,
        "duration_months": 1,
        "amount": 99.9,
        "payment_method": "alipay",
    }
    r = http.post(f"{base_url}/api/membership/orders", json=order_payload)
    assert r.status_code == 200, r.text
    data = r.json()
    assert data["success"] is True
    order_id = data["order_id"]

    # 查这个用户所有订单
    r = http.get(f"{base_url}/api/membership/orders/{user_id}")
    assert r.status_code == 200, r.text
    orders = r.json()
    assert isinstance(orders, list)
    assert len(orders) >= 1

    # 查最近一条
    r = http.get(f"{base_url}/api/membership/orders/{user_id}/latest")
    assert r.status_code == 200, r.text
    latest = r.json()
    assert latest["user_id"] == user_id

    # 查最近N条
    r = http.get(f"{base_url}/api/membership/orders/{user_id}/recent", params={"n": 2})
    assert r.status_code == 200, r.text
    recent_list = r.json()
    assert isinstance(recent_list, list)
