FROM python:3.5
ADD . /app
WORKDIR /app
#USER thedevilsvoice:thedevilsvoice

RUN pip install --upgrade pip
RUN pip install Cython
RUN pip install -r requirements/requirements.txt
CMD ["python", "bitsmasher/bitsmasher.py"]
