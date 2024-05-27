FROM mysql:5.7

# Copier les fichiers init.sql et init.sh dans les répertoires appropriés
COPY init.sql /docker-entrypoint-initdb.d/

# Afficher un message pendant la construction de l'image
RUN echo "Message de test : Construction de l'image en cours..."

# Définir les permissions sur les fichiers

RUN chown -R mysql:mysql /docker-entrypoint-initdb.d/ && \
    chmod 775 /docker-entrypoint-initdb.d/

