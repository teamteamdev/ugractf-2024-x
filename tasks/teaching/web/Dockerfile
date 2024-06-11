FROM python:alpine
RUN pip install quart uvicorn
WORKDIR /app
CMD ["uvicorn", "--uds", "/tmp/app.sock", "main:app"]
