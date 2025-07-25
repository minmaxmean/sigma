File of .siq format represents a SIGame questions package. It is actually a zip-archive with files containing questions and embedded multimedia content.

File names inside archive may be encoded in URI-format (for example, file _Снимок6.PNG_ can be stored as _%D0%A1%D0%BD%D0%B8%D0%BC%D0%BE%D0%BA6.PNG_). This encoding is applied for backward compatibility and it is not mandatory.

## Folder and file structure inside archive

* **content.xml**: main XML file which contains package description and questions text
* **Images**: folder for storing images
* **Audio**: folder for storing audio
* **Video**: folder for storing video
* **Html**: folder for storing HTML content

## content.xml file format

See [XSD schema](https://github.com/VladimirKhil/SI/blob/master/assets/siq_5.xsd) for full information.

Common file structure if the hierarchy of Package **(package)** - Round **(round)** - Theme **(theme)** - Question **(question)**, where each top-level element could contain arbitrary amount of low-level elements. Each item could contain some information (tag **info**) containing authors (**authors**), sources (**sources**), comments (**comments**) and showman comments (**showmanComments**) info. Authors and sources are inherited in hierarchy so, for example, if round does not contain authors information, package-level authors are considered to be the authors of this round.

Each author and source could be represented as common text or as a link to an entry in **global** authors and sources. A link starts with symbol **@** followed by item id.

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

Package entry could contain **global** tag having globally defined authors and sources (this is the replacement for **authors.xml** and **sources.xml** files of previous format).

Question entry contains type description (**type** attribute), parameters (**params**), right (**right**) and wrong (**wrong**) answers.

**script** entry could appear in question too but it is reserved for very complex scenarios in the future.

Question type is defined by a well-known or a custom name. Well-known question types are defined [here](https://github.com/VladimirKhil/SI/wiki/Question-types).

Question parameters have name and value. For each of well-known types there is a predefined set of parameters.

Each parameter has a type which could be equal to:

* simple (default): parameter has a string value
* content: parameter contains a set of content items (**item**)
* group: parameter contains a set of other parameters
* numberSet: parameter represents a range of numbers. The range is defined by minimum value (**minimum**), maximum value (**maximum**) and step value (**step**)

Each question contains a mandatory parameter called **question** having **content** type. This parameter defines a content to display when question is played.

Each content item is the minimal fragment of question play, having type, value and (optional) duration.

Content items of types **image**, **audio** and **video** could contain either an internal link to a file inside the package (then they have **isRef** attribute equal to **true**) or a link to an external file.

Additionally each content item has **placement** attribute which could be equal to:

* screen (default): this content is displayed on game screen
* replic: this content is provided as showman replic
* background: this content is played in the background

There is also **waitForFinish** attribute (true by default) which defines should the game wait for the content to finish playing before moving on or not. This attribute allow to display multiple content items on game screen.