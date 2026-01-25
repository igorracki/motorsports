import uvicorn
from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import Response
from ..services import SessionService, EventService
from ..providers import FastF1Provider

server = FastAPI(title="FastF1 Wrapper API")
server.add_middleware(
    CORSMiddleware,
    allow_origins=['*'],
    allow_methods=['*'],
    allow_headers=['*']
)


@server.get('/wrapper/race-result/{year}/{round}/{session_type}')
async def get_race_result(year: int, round: int, session_type: str):
    try:
        service = SessionService(FastF1Provider())
        result = service.get_race_result(year, round, session_type)
        return Response(content=result.json(), media_type="application/json")
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@server.get('/wrapper/events/{year}')
async def get_events(year: int):
    try:
        service = EventService(FastF1Provider())
        result = service.get_weekend_events(year)
        # rely on FastAPI's automatic JSON serialiation
        return result
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

if __name__ == '__main__':
    uvicorn.run(server, host='0.0.0.0', port=8080)
