$(function() {
  var video = document.getElementById("video");
  if (navigator.mediaDevices && navigator.mediaDevices.getUserMedia) {
    navigator.mediaDevices.getUserMedia({ video: true }).then(function(stream) {
      video.src = window.URL.createObjectURL(stream);
      video.play();
    });
  }

  var canvas = document.getElementById("canvas");
  var context = canvas.getContext("2d");
  var video = document.getElementById("video");

  // Trigger photo take
  var button = $("#snap");
  button.click(function() {
    button.addClass("loading");
    $(".info.message").hide();
    context.drawImage(video, 0, 0, 400, 225);
    var dataURL = canvas.toDataURL();
    $.ajax({
      type: "POST",
      url: "/webFaceID",
      data: {
        imgBase64: dataURL
      },
      success: function(resp) {
        const text =
          "standard: " +
          resp.tags.map(match => match.Tag) +
          "\n" +
          "custom:" +
          resp.custom_tags.map(match => match.Tag);
        $(".info.message")
          .text(text)
          .fadeIn();
        button
          .empty()
          .append($("<i>", { class: "camera icon" }))
          .addClass("teal")
          .removeClass("green");
          if (resp.custom_tags.length === 1) {
            window.location.replace("https://www.google.com/search?output=search&tbm=shop&q=" + resp.custom_tags[0].Tag);
          }
      },
      error: function() {
        $(".info.message")
          .text("Oops, something went wrong")
          .fadeIn();
      },
      complete: function() {
        button.removeClass("loading");
      }
    });
  });
});
