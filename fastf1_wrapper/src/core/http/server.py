import uvicorn
import logging
import sys
from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from ..services import F1Service
from ..providers import FastF1Provider

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[logging.StreamHandler(sys.stdout)]
)

server = FastAPI(title="FastF1 Wrapper API")
server.add_middleware(
    CORSMiddleware,
    allow_origins=['*'],
    allow_methods=['*'],
    allow_headers=['*']
)


@server.get('/wrapper/events/{year}')
async def get_events(year: int):
    try:
        service = F1Service(FastF1Provider())
        result = service.get_weekend_events(year)
        return result
    except Exception as e:
        logging.error(f"Error fetching events for year {year}: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@server.get('/wrapper/results/{year}/{round}/{session_type}')
async def get_results(year: int, round: int, session_type: str):
    try:
        service = F1Service(FastF1Provider())
        result = service.get_session_results(year, round, session_type)
        if result is None:
            return {"year": year, "round": round, "session_type": session_type, "results": []}
        return result
    except Exception as e:
        logging.error(f"Error fetching results for {year} round {round} {session_type}: {e}")
        raise HTTPException(status_code=500, detail=str(e))

if __name__ == '__main__':
    uvicorn.run(server, host='0.0.0.0', port=8080)
