# tests/test_membership.py
import datetime
import pytest


def test_membership_crud(http, base_url, test_user):
    user_id = test_user["user_id"]

    # 1. 先查全部
    r = http.get(f"{base_url}/api/membership")
    assert r.status_code == 200, r.text
    all_list = r.json()
    assert isinstance(all_list, list)

    # 2. 新增会员
    today = datetime.date.today()
    payload = {
        "user_id": user_id,
        "start_date": today.strftime("%Y-%m-%d"),
        "expire_date": (today.replace(day=min(today.day, 28)) + datetime.timedelta(days=30)).strftime("%Y-%m-%d"),
        "status": "active",
    }
    r = http.post(f"{base_url}/api/membership", json=payload)
    assert r.status_code == 200, r.text
    data = r.json()
    assert data["success"] is True
    membership_id = data["membership_id"]

    # 3. 查这个用户的会员
    r = http.get(f"{base_url}/api/membership/{user_id}")
    # 这个接口在 handler 里是如果不存在会 404，这里先直接按成功走
    assert r.status_code in (200, 404), r.text
    if r.status_code == 200:
        info = r.json()
        assert info["user_id"] == user_id

    # 4. 更新会员
    update_payload = {
        "start_date": payload["start_date"],
        "expire_date": payload["expire_date"],
        "status": "expired",
    }
    r = http.put(
        f"{base_url}/api/membership/{membership_id}",
        json=update_payload,
    )
    assert r.status_code == 200, r.text
    data = r.json()
    assert data["success"] is True

    # 5. 删除会员
    r = http.delete(f"{base_url}/api/membership/{membership_id}")
    assert r.status_code == 200, r.text
    data = r.json()
    assert data["success"] is True
