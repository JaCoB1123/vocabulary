package cmd

import (
	"encoding/json"
	"net/http"
	"os"
	"net/url"
	"net/http/httputil"

	"path"

	voc "github.com/JaCoB1123/vocabulary/internal/vocabulary"
	"github.com/JaCoB1123/web"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(webCmd)
}

type vocabularyServer struct {
	*web.Server
	*voc.Vocabulary
}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Serve vocabulary trainer via HTTP",
	Long:  "Used to access the vocabulary via a webbrowser",
	Run: func(cmd *cobra.Command, args []string) {
		server := &vocabularyServer{
			Server: web.NewServer(),
			Vocabulary: voc.MustVocabulary(*wordsfilename, *statsfilename),
		}

		remote, err := url.Parse("http://localhost:5000")
		if err != nil {
				panic(err)
		}
	
		proxy := httputil.NewSingleHostReverseProxy(remote)

		server.Server.Get("/api/stats", server.getStats)
		server.Server.Get("/api/learn", server.learn)
		server.Server.Post("/api/answer", server.answer)
		server.Server.Get("/(.*)", proxyHandler(proxy))


		http.ListenAndServe(":9876", server.Server)
	},
}
		
func proxyHandler(p *httputil.ReverseProxy) func(ctx *web.Context) {
	return func(ctx *web.Context) {
		p.ServeHTTP(ctx.ResponseWriter, ctx.Request)
	}
}

func (s *vocabularyServer) getStats(ctx *web.Context){
	stats := s.Vocabulary.GetVocabularyStats([]string{})
	jsonEncoder := json.NewEncoder(ctx.ResponseWriter)
	jsonEncoder.Encode(stats)
}

func (s *vocabularyServer) learn(ctx *web.Context){
	words := s.Vocabulary.GetSortedWords([]string{})
	jsonEncoder := json.NewEncoder(ctx.ResponseWriter)
	jsonEncoder.Encode(words[:10])
}

type answerRequest struct {
	Word string
	Success bool
}

func (s *vocabularyServer) answer(ctx *web.Context){
	req := answerRequest{}
	jsonDecoder := json.NewDecoder(ctx.Request.Body)
	jsonDecoder.Decode(&req)

	stats := s.Vocabulary.GetStats(req.Word)
	if req.Success {
		stats.CorrectAnswer()
	} else {
		stats.FalseAnswer()
	}

	s.Vocabulary.Save(*wordsfilename, *statsfilename)
}

var	STATIC_DIR = "web/dist"

func (s *vocabularyServer) static(ctx *web.Context, filepath string){
	if tryServeFile(ctx, filepath) {
		return
	}

	http.ServeFile(ctx.ResponseWriter, ctx.Request, path.Join(STATIC_DIR, "__app.html"))
}

func tryServeFile(ctx *web.Context, filepath string) bool {
	if filepath == "/" || filepath == "" {
		return false
	}

	filepath = path.Join(STATIC_DIR, filepath)
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	
	http.ServeFile(ctx.ResponseWriter, ctx.Request, filepath)
	return true
}