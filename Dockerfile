FROM centos
# 中文支持
ENV LANG en_US.UTF-8  
# 拷贝数据
WORKDIR /app
ADD ./task .
ADD etc etc
EXPOSE 8091
ENV ENV=production
CMD [ "./task","-init","-start" ]