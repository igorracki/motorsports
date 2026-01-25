FROM python:3.14-slim
WORKDIR /app
COPY pyproject.toml ./
RUN pip install -e .
COPY src/ ./src/
EXPOSE 8080
CMD ["python",  "-m", "src.core.http.server"]
