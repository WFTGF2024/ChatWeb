# tests/test_auth.py
import uuid
import pytest


def test_register_again(http, base_url):
    username = f"pytest_{uuid.uuid4().hex[:8]}"
    payload = {
        "username": username,
        "password": "Passw0rd!",
        "full_name": "Another User",
        "email": f"{username}@example.com",
        "phone_number": "13900000000",
        "security_question1": "q1",
        "security_answer1": "a1",
        "security_question2": "q2",
        "security_answer2": "a2",
    }
    r = http.post(f"{base_url}/api/auth/register", json=payload)
    assert r.status_code == 200, r.text
    data = r.json()
    assert data["success"] is True
    assert "user_id" in data


def test_login(http, base_url, test_user):
    r = http.post(
        f"{base_url}/api/auth/login",
        json={"username": test_user["username"], "password": test_user["password"]},
    )
    assert r.status_code == 200, r.text
    data = r.json()
    assert data["success"] is True
    assert "token" in data
    assert "expire_at" in data


def test_me(http, base_url, auth_headers):
    r = http.get(f"{base_url}/api/auth/me", headers=auth_headers)
    assert r.status_code == 200, r.text
    data = r.json()
    # 这里 handler 是直接把 profile 返回出来
    assert "username" in data
    assert "email" in data


def test_verify_security_and_reset(http, base_url, test_user):
    # 先验证密保
    verify_payload = {
        "username": test_user["username"],
        "security_answer1": "a1",
        "security_answer2": "a2",
    }
    r = http.post(f"{base_url}/api/auth/verify-security", json=verify_payload)
    assert r.status_code == 200, r.text
    data = r.json()
    assert data["success"] is True
    reset_token = data["reset_token"]

    # 再用 reset_token 改密码
    new_pass = "NewPassw0rd!"
    reset_payload = {
        "reset_token": reset_token,
        "new_password": new_pass,
    }
    r = http.post(f"{base_url}/api/auth/reset-password", json=reset_payload)
    assert r.status_code == 200, r.text
    data = r.json()
    assert data["success"] is True

    # 用新密码再登一次，确认真改了
    r = http.post(
        f"{base_url}/api/auth/login",
        json={"username": test_user["username"], "password": new_pass},
    )
    assert r.status_code == 200, r.text
    data = r.json()
    assert data["success"] is True
    assert "token" in data
