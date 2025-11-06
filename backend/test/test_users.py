# tests/test_users.py
import pytest


def test_update_user(http, base_url, test_user, auth_headers):
    user_id = test_user["user_id"]
    payload = {
        "full_name": "Pytest User Updated",
        "email": f"{test_user['username']}+updated@example.com",
        "phone_number": "13700000000",
    }
    r = http.put(
        f"{base_url}/api/users/{user_id}",
        json=payload,
        headers=auth_headers,
    )
    assert r.status_code == 200, r.text
    data = r.json()
    assert data["success"] is True


def test_delete_user(http, base_url, auth_headers):
    """
    为了不影响前面 fixture 里的主用户，这里临时再建一个用户然后删。
    """
    # 先建
    import uuid
    tmp_name = f"tmpdel_{uuid.uuid4().hex[:8]}"
    payload = {
        "username": tmp_name,
        "password": "Passw0rd!",
        "full_name": "To Be Deleted",
        "email": f"{tmp_name}@example.com",
        "phone_number": "13600000000",
        "security_question1": "q1",
        "security_answer1": "a1",
        "security_question2": "q2",
        "security_answer2": "a2",
    }
    r = http.post(f"{base_url}/api/auth/register", json=payload)
    assert r.status_code == 200, r.text
    user_id = r.json()["user_id"]

    # 再删。注意：你的 handler 是从 token 里取 user_id，所以这里要用这个新用户的 token
    # 先登录拿 token
    r = http.post(
        f"{base_url}/api/auth/login",
        json={"username": tmp_name, "password": "Passw0rd!"},
    )
    assert r.status_code == 200, r.text
    tmp_token = r.json()["token"]
    headers = {"Authorization": f"Bearer {tmp_token}"}

    r = http.delete(f"{base_url}/api/users/{user_id}", headers=headers)
    assert r.status_code == 200, r.text
    data = r.json()
    assert data["success"] is True
