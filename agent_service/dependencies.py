from fastapi import Depends, Request

def get_db_client(request: Request):
    if not hasattr(request.app.state, "mongodb_client"):
        raise HTTPException(status_code=500, detail="Database not initialized")
    return request.app.state.mongodb_client

