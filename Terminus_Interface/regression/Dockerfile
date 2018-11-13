FROM continuumio/miniconda3
RUN conda update -y conda && conda update -y python && \
    pip install --upgrade pip
COPY requirements.txt /
RUN pip install --no-cache-dir -r requirements.txt
COPY src/ /app
EXPOSE 9002
WORKDIR /app
CMD ["python","-u", "main.py"]
