$(document).ready(function() {
  var result = '{{ .Result | safehtml }}';
  if (result != "ERROR"){
    $("#result").text(JSON.stringify(JSON.parse(result), null, "\t"));
    $("#info").removeClass("hidden").addClass("visible");
  } else {
    console.log(result);
    $("#warning").text(result).removeClass("hidden").addClass("visible");
  }
});
var request = function(data, callback, onerror) {
  $.ajax({
    type:          'POST',
    dataType:      'json',
    contentType:   'application/json',
    scriptCharset: 'utf-8',
    data:          JSON.stringify(data),
    url:           App.url
  })
  .done(function(res) {
    callback(res);
  })
  .fail(function(e) {
    onerror(e);
  });
};
