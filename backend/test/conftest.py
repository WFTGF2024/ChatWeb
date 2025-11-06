# tests/conftest.py
import os
import uuid
import random
import pytest
import requests


@pytest.fixture(scope="session")
def base_url():
    # 你也可以在命令行里 `BACKEND_BASE_URL=http://127.0.0.1:8080 pytest`
    return os.getenv("BACKEND_BASE_URL", "http://localhost:5000")


@pytest.fixture(scope="session")
def http(base_url):
    s = requests.Session()
    # 用一个最“便宜”的接口测一下服务在不在
    try:
        resp = s.post(f"{base_url}/web/chunk", json={})
    except Exception as e:
        pytest.skip(f"backend is not reachable at {base_url}: {e}")
    else:
        # 200/404/405 都算你服务活着
        if resp.status_code not in (200, 404, 405):
            pytest.skip(f"unexpected status from backend: {resp.status_code}")
    return s


def _random_phone():
    # 简单搞一个 11 位手机号
    return "138" + "".join(str(random.randint(0, 9)) for _ in range(8))


@pytest.fixture(scope="session")
def test_user(http, base_url):
    """注册 + 登录，返回 user_id 和 token"""
    username = f"pytest_{uuid.uuid4().hex[:8]}"
    password = "Passw0rd!"
    payload = {
        "username": username,
        "password": password,
        "full_name": "Pytest User",
        "email": f"{username}@example.com",
        "phone_number": _random_phone(),
        "security_question1": "q1",
        "security_answer1": "a1",
        "security_question2": "q2",
        "security_answer2": "a2",
    }

    # 注册
    r = http.post(f"{base_url}/api/auth/register", json=payload)
    assert r.status_code == 200, r.text
    data = r.json()
    assert data.get("success") is True
    user_id = data.get("user_id")
    assert user_id, "register should return user_id"

    # 登录
    r = http.post(
        f"{base_url}/api/auth/login",
        json={"username": username, "password": password},
    )
    assert r.status_code == 200, r.text
    data = r.json()
    token = data.get("token")
    assert token, "login should return token"

    return {
        "user_id": user_id,
        "username": username,
        "password": password,
        "token": token,
    }


@pytest.fixture
def auth_headers(test_user):
    return {"Authorization": f"Bearer {test_user['token']}"}
