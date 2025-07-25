File of .siq format represents a SIGame questions package. It is actually a zip-archive with files containing questions and embedded multimedia content.

File names inside archive may be encoded in URI-format (for example, file _Снимок6.PNG_ can be stored _%D0%A1%D0%BD%D0%B8%D0%BC%D0%BE%D0%BA6.PNG_). This encoding is applied for backward compatibility and it is not mandatory.

## Folder and file structure inside archive

* **content.xml**: main XML file which contains package description and questions text
* **[Content_Types].xml**: auxiliary file which exists only for backward compatibility. There is no nee to use it
* **Texts**: folder which contains information about global package authors and sources:
  * **authors.xml**: XML file with global authors info
  * **sources.xml**: XML file with global sources info
* **Images**: folder for storing images
* **Audio**: folder for storing audio
* **Video**: folder for storing video

## content.xml file format

See [XSD schema](https://github.com/VladimirKhil/SI/blob/master/assets/ygpackage3.1.xsd) for full information.

Common file structure if the hierarchy of Package **(package)** - Round **(round)** - Theme **(theme)** - Question **(question)**, where each top-level element could contain arbitrary amount of low-level elements. Each item could contain some information (tag **info**) containing authors (**authors**), sources (**sources**) and comments (**comments**) info. Authors and sources are inherited in hierarchy so, for example, if round does not contain authors information, package-level authors are considered to be the authors of this round.

Each author and source could be represented as common text or as a link to an entry in files authors.xml/sources.xml. A link starts with symbol **@** followed by item id.

Example: `<author>@ae9f7eb2-6091-4b34-97a1-0f74ad193d57</author>`

Moreover, a source could contain specification after the link id (for example, page number). Specification is denoted by symbol **#** placed after source id. After **#** goes specification text.

Example: `<source>@7ab08cfa-7f68-4fd4-a400-f4ac26b33a9d#с.256</source>`

Package item contains following additional attributes:

* **id**: unique package id
* **name**: package name
* **version**: package schema version
* **restriction**: package restrictions
* **date**: package creation date
* **publisher**: package publisher
* **difficulty**: package difficulty (from 1 to 10)
* **logo**: package logo link (symbol **@** and then a file name from Images folder)
* **language**: package language

Also package entry could contain tags entry (**tags**), where each element (**tag**) describe some package theme ("Cinema", "Games", "Books" etc).

Question entry contains type description (**type**), scenario (**scenario**), right (**right**) and wrong (**wrong**) answers.

Question type is defined by a well-known or a custom name. It has a set of parameters specific for each type. Well-known question types are defined [here](http://vladimirkhil.com/si/qtypes).

Question scenario contains one or several atoms (**atom**). An atom is the minimal fragment of question play, which has a type and (optional) duration. There is also an atom of the special type **marker** which denotes that all atoms after it belong to answer and not the question body. With such a market it is possible to escribe complex answers.

Atoms of types **image**, **voice** and **video** could contain either internal link to a file inside a package (such links texts start with **@** followed by a file name) or the link to an external file.

## authors.xml file format

File contains entries of type **Author** with fields:

* **id**: unique author identifier (Guid)
* **Name**: author name
* **SecondName**: author second name
* **Surname**: author surname
* **Country**: author country
* **City**: author city

## sources.xml file format

File contains entries of type **Source** with fields:

* **id**: unique source identifier (Guid)
* **Author**: source author
* **Title**: source title
* **Year**: source publish (record) year
* **Publish**: source publisher (company)
* **City**: source publisher city