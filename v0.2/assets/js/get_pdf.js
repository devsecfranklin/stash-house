function descargar_pdf() {
  var pdf = new jsPDF("p", "pt");

  var specialElementHandlers = {
    "#editor": function (element, renderer) {
      return true;
    },
  };

  pdf.setProperties({
    title: "Franklin Diaz Resume",
    subject:
      "resume of a DevSecOps/security infrastructure engineer seeking remote opportunities.",
    author: "Franklin Diaz",
    keywords: "programatically generated, javascript, web 2.0, ajax",
    creator: "Franklin Diaz, 2022",
  });

  var options = {
    pagesplit: true,
  };

  pdf.addHTML(document.body, options, function () {
    pdf.save("franklin-diaz-resume.pdf");
  });
}
function descargar_pdf() {
  var pdf = new jsPDF("p", "pt");

  var specialElementHandlers = {
    "#editor": function (element, renderer) {
      return true;
    },
  };

  pdf.setProperties({
    title: "Franklin Diaz Resume",
    subject:
      "resume of a DevSecOps/security infrastructure engineer seeking remote opportunities.",
    author: "Franklin Diaz",
    keywords: "programatically generated, javascript, web 2.0, ajax",
    creator: "Franklin Diaz, 2022",
  });

  var options = {
    pagesplit: true,
  };

  pdf.addHTML(document.body, options, function () {
    pdf.save("franklin-diaz-resume.pdf");
  });
}
