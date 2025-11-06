# tests/test_chat_web.py
def test_web_ingest(http, base_url):
    r = http.post(f"{base_url}/web/ingest", json={"url": "https://example.com"})
    assert r.status_code == 200, r.text
    data = r.json()
    assert data["success"] is True


def test_web_chunk(http, base_url):
    r = http.post(f"{base_url}/web/chunk", json={"content": "hello"})
    assert r.status_code == 200, r.text
    data = r.json()
    assert data["success"] is True


def test_chat_history(http, base_url):
    r = http.get(f"{base_url}/chat/history")
    assert r.status_code == 200, r.text
    data = r.json()
    assert data["success"] is True


def test_chat_file_bridge(http, base_url):
    r = http.post(f"{base_url}/chat/file-bridge", json={"file_id": "dummy"})
    assert r.status_code == 200, r.text
    data = r.json()
    assert data["success"] is True
