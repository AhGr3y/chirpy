FROM debian:stable-slim
COPY chirpy /bin/chirpy
ENV JWT_SECRET=Jeoo7719AyAkDHMbpmZzX8MzULTtTHxibhtx6Dqp3sRTOqzoN/icUMYJhqF8r0WioxFInMyY8lGZEfpF90sWeQ==
ENV POLKA_KEY=f271c81ff7084ee5b99a5091b42d486e
CMD ["/bin/chirpy"]