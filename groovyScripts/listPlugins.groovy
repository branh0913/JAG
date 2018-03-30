import jenkins.model.*
import hudson.model.*


Jenkins.instance.pluginManager.plugins.each{
    plugin -> println("${plugin.getDisplayName()}")
}