
# RubberyConf Software Configuration DSL

## Abstract

RuberyConf is a DSL language for general-purpose software configuration management focusing on business. It describes an abstraction layer to standardize how software can be configured for business features. It supports concurrent managing configurations based on a set of rules. It also enables non-engineers to edit and change business configurations dynamically in a standard. RubberyConf is extensible, and it allows contributors to extend its features.  

### Tags

software configuration, configuration as code, CaC, feature toggles, feature flags

## Introduction

'Application configuration' is a concept usually realized differently by distinct stakeholders in a project. Although the business requirements are the source of truth for a project, these requirements can be realized differently in each layer till they arrive at the engineering team. There are several methodologies in the industry like CMMI-Dev[1] or Software Configuration Management [2], they describe several processes to build a software item, but it doesn’t cover how a feature should be described. Also, in today's world, mobile apps have radically different needs when compared to backend needs. As software engineers, we should find a way to enable our stakeholders (mainly product owners) to enable or disable configurations without redeploying app or services, even better without involving engineers in testing those configurations in production when they’ve been updated [3]. Unfortunately, software configuration isn’t well described in the industry and there are a few documents in the tech literature. Often, this topic is related to DevOps practices too, but it goes in system configuration instead of business. Most of those documents are focused on tech challenges instead of configuration needs. Some of them describe Feature Toggle systems[4,5]. But these systems have been built ad-hoc based on concrete needs. This document describes a high-level DSL language to describe general-purpose software configuration. Low-level details to implement a system such as database connections and configurations, as well as third-party services, are out of the scope of this DSL.   

## Use Case

Assume that we would like to enable a Product Owner (PO) to update software configuration either by enabling or disabling features. Often, when a software development team enables a feature in production, they want to manage this rolling out process carefully. Let’s assume we have a software part of which was updated recently. Let's further assume, this piece of software has algorithm A and algorithm B. Now we want to enable a PO to switch from algorithm A to algorithm B by themselves (no engineer activity involved toggling this logic in production).

An approach to solve this problem might be using environment variables (EV), but it requires shell script knowledge and a good understanding of operative systems to change values. An alternative might be storing this configuration in an RDBMS database (MySQL, Oracle or Microsoft SQLServer). Still, it requires familiarity with SQL commands, and often production environments are locked to manual changes. Some companies are using Feature Toggle, but they are created ad-hoc for their needs. Indeed, it isn’t reusable across the industry as a standard.

Furthermore, PO might face more complex requirements about targeting the feature switch. For example, they might need to enable algorithm B for a concrete subset of users based on some rules. Meantime, A version might still serve value to other users.  

Finally, engineers must develop both A and B versions of their software. However, if A or B is enabled, it must be transparent for engineers. 

## RuberyConf as a configuration language

We’ve created a DSL to describe an abstraction layer to change software configuration with a set of predefined actions. We’ve proposed to use YAML as a base syntax to describe those actions. 

RubreryConf is a simple set of fields, although extensible, in yml format. A simple configuration must describe our use case like:

```yml
name: Algorithm
default:  
  value: 
    data: false
    type: boolean
```

Given this configuration, let's assume that algorithm A will be executed when data is 'false' and algorithm B will be executed when it is true in the implementation layer. With this simple change, production software will stop serving code based on A to resolve with the B version.

Type field describes what type of data must be expected. It can be boolean, number, json, and file. So, the data field must be aligned with it.

We want to enable the B version for a concrete scenario, so RuberyConf allows adding alternative configurations based on some circumstances, for example:

```yml
configurations:
  - config: 
    rulesBehaviour: "AND"
    rules:
      platform:
        - android
      version:
       min: “1.0.0”
       max: “10.0.0”
      country:
        - UK
    rollout:
      enabledForOnly: 0.10
    value: true

```

In the above example, there is one config enabled with true value. All these rules must be satisfied to rewrite the default value (false) to true. RulesBehaviour field manages whether all rules must be satisfied (AND) or with one of them is enough (OR) to rewrite the default value. The rollout field describes how this rule must be applied; in this case, it means only 10% of users will get true. For example, 10% of users with android devices and version 5.0.1 in the United Kingdom country will get true, other users will see false as the default value. 

All fields are extensible with new values and rules. 

Engineers must implement those fields manually in their software. Fields must be checked against values coming from software infrastructure.  

### RubberyConf syntax

#### Meta field

Adds meta-information about this configuration. Meta has two fields: description and tags (both Optional).

```yml
meta:
  description: 'example 1'
  tags:
    - Foo
    - Woo
```

#### Default field

Default (mandatory) is the main description of this configuration. It has these fields: ttl (optional) and value.  Ttl sets how long this configuration will be in memory. Value stores the default configuration and has 2 fields: data and type. Data is the real value to be sent to the customer and type describes what data type is, both mandatory. Example

```yml
default:  
  ttl: 120s
  value: 
    data: "{value:0}"
    type: jsonFormat
```

#### Configurations field

The default value can be overwritten by different use cases or configs. Each config is an entry of field configurations. Field configurations are an array of config.

#### Config field

Config field represents a set of rules and a value. When the set of rules is satisfied then default value field is overwritten by this config value field. Rules is an array of rule fields (described below). Rules can work as logical operators such as ‘and’ or ‘or’. When ‘and’ is set all rules must be true, otherwise ‘or’ means with at least one being true is enough to overwrite default value by this config value.  Example: 

```yml
configurations:
  - config: 
    rulesBehaviour: "OR"
    rules:
      - rule 1
      - rule 2
    rollout:
      enabledForOnly: 1.0
    value: "{value:1}"
```

#### Rule field

Rule field is an abstraction of the code that will be implemented in the software to be configured. List of rules are: environment, querystring, header, platform, version, country, city, userID, userGroup and timer. This list is extensible.  Example:

```yml
    rules:
      - querystring:
        value: 
          - 23342
        key: "userId"
      - header:
        value:
          - cyan
        key: "color"
      -  environment:
          - production
          - stage
```

#### Rollout field

It has one field: enabledForOnly, it means the share of users will get a value. Example: 

```yml
    rollout: 
      enabledForOnly: 0.10

```

#### List of examples

A list of examples from simple to high complexity is in RubberyConf language repository [6].

## Conclusion

With this abstraction layer, all general-purpose software can be configured in a standard manner.  Even non-engineers can change configuration without technical knowledge. Just changing a configuration in the yml file can apply business changes dynamically. RubberyConf isn’t an ABTesting tool, although it might be used for testing purposes. RubberyConf follows the ‘Config as Code’ (CaC) pattern. So, we can keep these configuration files in any CVS and use their features to track and revert changes, authorization and authentication and so on.

## Acknowledgments 

I would like to thank reviewers [Erol Aran](https://github.com/erolaran) and [Gastón Fournier](https://github.com/gastonfournier) for their valuable suggestions.


## References

* CaC https://dzone.com/articles/infrastructure-versus-config-as-code 
* PO https://www.scrum.org/resources/what-is-a-product-owner 
* DSL https://en.wikipedia.org/wiki/Domain-specific_language 
* CVS https://en.wikipedia.org/wiki/Concurrent_Versions_System 
* API https://en.wikipedia.org/wiki/Web_API 
* RDBMS https://en.wikipedia.org/wiki/Relational_database#RDBMS
* YAML, YML file https://es.wikipedia.org/wiki/YAML
* DevOps https://engineering.atspotify.com/2013/05/17/devops-management/
* [1] CMMI for development, v1.3 ​​https://resources.sei.cmu.edu/asset_files/TechnicalReport/2010_005_001_15287.pdf 
* [2] 828-2012 IEEE Standard for Configuration Management in Systems and Software Engineering. 2012
* [3] Chunqiang Tang, Thawan Kooburat, Pradeep Venkatachalam, Holistic Configuration Management @ Facebook 2016, pp. 7-8 https://research.fb.com/wp-content/uploads/2016/11/holistic-configuration-management-at-facebook.pdf 
* [4]Feature toggle, https://martinfowler.com/articles/feature-toggles.html 
* [5]Feature toggles at Netflix, https://netflixtechblog.com/preparing-the-netflix-api-for-deployment-786d8f58090d 
* [6] https://github.com/rubberyconf/language  



