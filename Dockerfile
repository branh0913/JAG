FROM jenkins/jenkins:lts

#COPY groovyScripts/installPlugins.groovy /usr/share/jenkins/ref/init.groovy.d/
#COPY groovyScripts/default-user.groovy /usr/share/jenkins/ref/init.groovy.d/
ENV JENKINS_OPTS --httpPort=3000
ENV JAVA_OPTS -Djenkins.install.runSetupWizard=false
EXPOSE 3000