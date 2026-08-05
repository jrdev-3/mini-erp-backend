package contact

import (
	"errors"
	"regexp"
	"strings"
)

// ContactRequest representa o payload enviado no corpo da requisição para criar ou atualizar um contato.
type ContactRequest struct {
	IsCustomer      bool    `json:"is_customer"`
	IsSupplier      bool    `json:"is_supplier"`
	IsActive        *bool   `json:"is_active"` // Usado como ponteiro para poder validar se foi enviado no JSON
	TipoPessoa      string  `json:"tipo_pessoa"` // PF ou PJ
	NomeRazaoSocial string  `json:"nome_razao_social"`
	Documento       string  `json:"documento"` // CPF ou CNPJ (será limpo para conter apenas números)
	Email           *string `json:"email"`
	Telefone        *string `json:"telefone"`
	Rua             *string `json:"rua"`
	Numero          *string `json:"numero"`
	Bairro          *string `json:"bairro"`
	Cidade          *string `json:"cidade"`
	Estado          *string `json:"estado"`
	CEP             *string `json:"cep"`
}

var (
	rxEmail = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,4}$`)
)

// Validate executa a validação sintática e de formatos dos dados de entrada do contato.
func (r *ContactRequest) Validate() error {
	r.NomeRazaoSocial = strings.TrimSpace(r.NomeRazaoSocial)
	if r.NomeRazaoSocial == "" {
		return errors.New("nome/razão social é obrigatório")
	}

	r.TipoPessoa = strings.ToUpper(strings.TrimSpace(r.TipoPessoa))
	if r.TipoPessoa != "PF" && r.TipoPessoa != "PJ" {
		return errors.New("tipo de pessoa inválido (deve ser PF ou PJ)")
	}

	// Limpar documento deixando apenas números
	r.Documento = cleanNumericString(r.Documento)
	if r.Documento == "" {
		return errors.New("documento (CPF/CNPJ) é obrigatório")
	}

	if r.TipoPessoa == "PF" {
		if len(r.Documento) != 11 {
			return errors.New("CPF deve conter exatamente 11 dígitos numéricos")
		}
	} else {
		if len(r.Documento) != 14 {
			return errors.New("CNPJ deve conter exatamente 14 dígitos numéricos")
		}
	}

	if r.Email != nil && *r.Email != "" {
		emailClean := strings.ToLower(strings.TrimSpace(*r.Email))
		r.Email = &emailClean
		if !rxEmail.MatchString(*r.Email) {
			return errors.New("formato de e-mail inválido")
		}
	}

	if r.CEP != nil && *r.CEP != "" {
		cepClean := cleanNumericString(*r.CEP)
		r.CEP = &cepClean
		if len(*r.CEP) != 8 {
			return errors.New("CEP deve conter exatamente 8 dígitos numéricos")
		}
	}

	if r.Estado != nil && *r.Estado != "" {
		estadoClean := strings.ToUpper(strings.TrimSpace(*r.Estado))
		r.Estado = &estadoClean
		if len(*r.Estado) != 2 {
			return errors.New("estado deve conter exatamente 2 caracteres (UF)")
		}
	}

	return nil
}

// cleanNumericString remove qualquer caractere que não seja um dígito numérico (0-9)
func cleanNumericString(s string) string {
	var sb strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}
