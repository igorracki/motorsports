import uvicorn
import logging
import sys
from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from ..services import F1Service
from ..providers import FastF1Provider

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[logging.StreamHandler(sys.stdout)]
)

logger = logging.getLogger(__name__)
server = FastAPI(title="FastF1 Wrapper API")
server.add_middleware(
    CORSMiddleware,
    allow_origins=['*'],
    allow_methods=['*'],
    allow_headers=['*']
)


@server.get('/wrapper/events/{year}')
async def get_events(year: int):
    logger.info(f"Entry: get_events(year={year})")
    try:
        service = F1Service(FastF1Provider())
        results = service.get_weekend_events(year)
        logger.info(f"Exit: get_events(year={year}) - Found {len(results)} events")
        return results
    except Exception:
        logger.exception(f"Error fetching events for year {year}")
        raise HTTPException(status_code=500, detail="Internal server error")


@server.get('/wrapper/results/{year}/{round}/{session_type}')
async def get_results(year: int, round: int, session_type: str):
    logger.info(f"Entry: get_results(year={year}, round={round}, session={session_type})")
    try:
        service = F1Service(FastF1Provider())
        result = service.get_session_results(year, round, session_type)
        if result is None:
            logger.info(f"Exit: get_results(year={year}, round={round}, session={session_type}) - No results found")
            return {"year": year, "round": round, "session_type": session_type, "results": []}
        
        count = len(result.results) if result.results else 0
        logger.info(f"Exit: get_results(year={year}, round={round}, session={session_type}) - Found {count} results")
        return result
    except Exception:
        logger.exception(f"Error fetching results for {year} round {round} {session_type}")
        raise HTTPException(status_code=500, detail="Internal server error")


@server.get('/wrapper/circuits/{year}/{round}')
async def get_circuit(year: int, round: int):
    logger.info(f"Entry: get_circuit(year={year}, round={round})")
    try:
        service = F1Service(FastF1Provider())
        result = service.get_circuit_data(year, round)
        if result is None:
            logger.warning(f"Exit: get_circuit(year={year}, round={round}) - Not found")
            raise HTTPException(status_code=404, detail="Circuit not found")
        
        logger.info(f"Exit: get_circuit(year={year}, round={round}) - Success")
        return result
    except HTTPException as e:
        raise e
    except Exception:
        logger.exception(f"Error fetching circuit for {year} round {round}")
        raise HTTPException(status_code=500, detail="Internal server error")


if __name__ == '__main__':
    uvicorn.run(server, host='0.0.0.0', port=8080)
