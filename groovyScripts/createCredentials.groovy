import com.cloudbees.plugins.credentials.impl.*
import com.cloudbees.plugins.credentials.*
import com.cloudbees.plugins.credentials.domains.*

Credentials c = (Credentials) new UsernamePasswordCredentialsImpl(CredentialsScope.GLOBAL,"{{.CredentialID}}", "description", "{{.ServiceAccount}}", "{{.Servicepass}}")
SystemCredentialsProvider.getInstance().getStore().addCredentials(Domain.global(), c)
println("Credentials have been created")